package inspector

import (
	"github.com/florianl/go-nfqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (i *inspector) handler(a nfqueue.Attribute) int {
	id := *a.PacketID
	payload := a.Payload

	packet := gopacket.NewPacket(
		*payload,
		layers.LayerTypeIPv4,
		gopacket.Default,
	)

	var dst string
	if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		ipv4, ok := ipLayer.(*layers.IPv4)
		if !ok {
			i.logger.Info().Msg("type assertion to ipv4 fail")

			if err := i.queue.SetVerdict(id, nfqueue.NfAccept); err != nil {
				i.logger.Warn().Err(err).Msg("set accept verdict for non ipv4 packet fail")
			}

			return 0
		}

		dst = ipv4.DstIP.String()

		if _, blocked := i.blockedAddresses.Load(dst); blocked {
			if err := i.queue.SetVerdict(id, nfqueue.NfDrop); err != nil {
				i.logger.Warn().Err(err).Msg("set drop verdict for packet fail")
				return 0
			}
		}
	}

	if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp, ok := tcpLayer.(*layers.TCP)
		if !ok {
			i.logger.Info().Msg("type assertion to tcp fail")

			if err := i.queue.SetVerdict(id, nfqueue.NfAccept); err != nil {
				i.logger.Warn().Err(err).Msg("set accept verdict for non tcp packet fail")
			}

			return 0
		}

		for _, detector := range i.detectors {
			check, verdict, err := detector.Detect(tcp.Payload)
			if err != nil {
				i.logger.Warn().Err(err).Msg("deep packet inspection fail")

				if err := i.queue.SetVerdict(id, nfqueue.NfAccept); err != nil {
					i.logger.Warn().Err(err).Msg("set accept verdict for proceed tcp packet fail")
					return 0
				}
			}

			if *check.ValuePtr() {
				i.blockedAddresses.Store(dst, verdict)

				if err := i.queue.SetVerdict(id, nfqueue.NfDrop); err != nil {
					i.logger.Warn().Err(err).Msg("set drop verdict for obfs packet fail")
					return 0
				}
			}
		}

		if err := i.queue.SetVerdict(id, nfqueue.NfAccept); err != nil {
			i.logger.Warn().Err(err).Msg("set accept verdict for clear packet fail")
			return 0
		}
	}

	return 0
}
