# staged_robot

Go source code for the **Staged Robot** LinMot Z-axis control system. This repository contains two Go modules:

| Directory                          | Module                                       | Purpose                                                                         |
| ---------------------------------- | -------------------------------------------- | ------------------------------------------------------------------------------- |
| [`linmot_client/`](linmot_client/) | `github.com/Smart-Vision-Works/staged_robot` | LinUDP V2 protocol library for LinMot C1250 drives                              |
| [`stage_primer/`](stage_primer/)   | `primer`                                     | REST/gRPC service that uses linmot_client to control staged robot Z-axis motion |

## How They Fit Together

```
                     ┌──────────────────────────────────────┐
                     │  TensorPro (tater_spider_sai)        │
                     │  Sends pick commands + settings       │
                     └──────────────┬───────────────────────┘
                                    │ gRPC :50051
                                    ▼
┌───────────────────────────────────────────────────────────────┐
│  stage_primer/                                                │
│                                                               │
│  server/grpc_server.go     ← DeployCommandTable, Jog, Setup  │
│         │                                                     │
│         ▼                                                     │
│  linmot/command_table.go   ← builds YAML template, binds     │
│  linmot/setup.go              Z-distance, speed, pick time   │
│  linmot/jog.go             ← calls into linmot_client ──────────┐
│  linmot/faults.go                                             │  │
│  linmot/vacuum.go                                             │  │
└───────────────────────────────────────────────────────────────┘  │
                                                                   │
┌──────────────────────────────────────────────────────────────────┘
│
│  linmot_client/
│
│  client/pool.go            ← SharedUDPTransport, one socket
│  client/client.go          ← SetCommandTable, ReadRAM, WriteRAMAndROM
│  client/rtc/               ← RTC command encoding
│  client/control_word/      ← EnableDrive, Home, AcknowledgeError
│  protocol/                 ← raw LinUDP V2 packet encode/decode
│  transport/                ← UDP I/O
│         │
│         ▼
│  LinMot C1250 drives (UDP port 49360)
```

## Where stage_primer Calls linmot_client

See [`stage_primer/README.md`](stage_primer/README.md) for the full callsite map with file and line references.

## Reference Documentation

Official LinMot hardware and protocol docs (converted from PDF to Markdown) are in [`linmot_client/reference/`](linmot_client/reference/).
