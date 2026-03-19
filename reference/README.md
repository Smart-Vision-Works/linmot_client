# LinMot Reference Documentation

This directory contains official LinMot documentation (converted to Markdown from NTI AG PDFs) and source code references for LinMot integration.

---

## Documents

### Hardware

| Document | Source PDF | Description |
|----------|------------|-------------|
| [C1250_MI_Installation_Guide.md](C1250_MI_Installation_Guide.md) | [C1250_MI_Installation_Guide.pdf](C1250_MI_Installation_Guide.pdf) | C1250-MI hardware installation guide. Connector pinouts (X1–X19), DIP switch S1/S2 address table, system LED blink codes, safety requirements, ordering info. |

### Protocol & Software

| Document | Source PDF | Description |
|----------|------------|-------------|
| [LinUDP_V2.md](LinUDP_V2.md) | [LinUDP_V2.pdf](LinUDP_V2.pdf) | LinUDP V2 protocol manual (v1.9). Default IP addressing via DIP switches, UDP packet format (request + response), UPID parameter list, monitoring channels, master configuration modes. |
| [LinUDP_V2_DLL.md](LinUDP_V2_DLL.md) | [LinUDP_V2_DLL.pdf](LinUDP_V2_DLL.pdf) | LinUDP DLL integration guide (v2.1.1). ACI class: all methods and functions, state machine handling, motion commands, command table access, drive parameter access, load/save configuration. **Appendix I** (commissioning) is the authoritative source for the "Static by IP Configuration" parameter path. |
| [LinMotTalk_Manual.md](LinMotTalk_Manual.md) | [LinMotTalk_Manual.pdf](LinMotTalk_Manual.pdf) | LinMot-Talk v6 software manual. Ethernet scan procedure, RS232 connection by drive type, UPID parameter editing, Save to ROM, and **§4.1 Factory Reset** — the authoritative 6-step hardware reset procedure for C1250 drives. |
| [LinMot_MotionCtrl.md](LinMot_MotionCtrl.md) | [LinMot_MotionCtrl.pdf](LinMot_MotionCtrl.pdf) | Motion Control Software user manual. Easy Steps (auto start, auto home, triggered command table), command table format, position controller settings (Set Soft vs. Set Stiff), homing configuration. |

### Source Code Reference

| File | Description |
|------|-------------|
| [decompiled_linudp_csharp_lib.cs](decompiled_linudp_csharp_lib.cs) | Decompiled official LinUDP C# library. Reference for understanding ACI function signatures, packet construction, and behavior that the Go implementation mirrors. |

---

## Key Reference Points

| Topic | Where to Look |
|-------|--------------|
| DIP switch bit layout (S1/S2) | `C1250_MI_Installation_Guide.md` §9.9 |
| LED blink codes (all states) | `C1250_MI_Installation_Guide.md` §10 |
| Factory reset procedure (6-step, hardware) | `LinMotTalk_Manual.md` §4.1 |
| "Waiting for Defaulting Parameters" LED state | `C1250_MI_Installation_Guide.md` §10 — ERROR+WARN alternating ~4 Hz |
| "Defaulting Parameters Done" LED state | `C1250_MI_Installation_Guide.md` §10 — WARN+EN together ~2 Hz |
| Default IP addressing (192.168.1.xxx via DIP switches) | `LinUDP_V2.md` §3.2 |
| Static IP configuration ("Static by IP Configuration") | `LinUDP_V2_DLL.md` Appendix I |
| LinUDP V2 packet format | `LinUDP_V2.md` §4 |
| LinUDP V2 UPID/parameter list | `LinUDP_V2.md` §6 |
| Monitoring channel UPIDs (0x20A8–0x20AB) | `LinUDP_V2.md` §5 |
| Commissioning the LinUDP interface in LinMot-Talk | `LinUDP_V2_DLL.md` Appendix I |
| RS232 connection on C1250 (X19, RJ45) | `C1250_MI_Installation_Guide.md` §9.7 |
| Required RS232 cable and adapter art. numbers | `C1250_MI_Installation_Guide.md` §11 (0150-2143, 0150-2473) |
| Easy Steps auto start / auto home | `LinMot_MotionCtrl.md` |
| Position controller: Set Soft vs. Set Stiff | `LinMot_MotionCtrl.md` |
| 14 operational ROM parameters written on commissioning | `LinUDP_V2_DLL.md` Appendix I |

---

## Drive Configuration Reference

Reference values for commissioning or recovering these drives:

| Parameter | LinMot 0 | LinMot 1 |
|-----------|----------|----------|
| Static IP | `10.8.7.232` | `10.8.7.234` |
| Subnet Mask | `255.255.248.0` | `255.255.248.0` |
| Default Gateway | `10.8.0.1` | `10.8.0.1` |
| LinUDP Drive Port | `49360` | `49360` |
| DIP switch (discovery mode) | `0xE8` (S1=0xE, S2=0x8) | `0xEA` (S1=0xE, S2=0xA) |
| Discovery IP (DIP-switch mode) | `192.168.1.232` | `192.168.1.234` |
| Actuator | `0150-4066_V3S7_20250612_DM01-23x80F-HP-R-160_MS11.adp` | same |
| Firmware interface | LinUDP + Easy Steps | same |
