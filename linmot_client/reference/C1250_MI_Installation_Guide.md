# C1250-MI Installation Guide — LinMot Servo Drive

**Document:** 0185-1175-E_1V11_IG_Drives_C1250-MI  
**Product:** C1250-MI Multi Interface Servo Drive  
**Manufacturer:** NTI AG / LinMot, Bodenaeckerstrasse 2, CH-8957 Spreitenbach, Switzerland  
**Latest version:** http://www.linmot.com

---

## 1 General Information

### 1.1 Introduction

This manual includes instructions for the assembly, installation, maintenance, transport, and storage of the servo drives. The document is intended for electricians, mechanics, service technicians, and warehouse staff. Read this manual before using the product and always observe the general safety instructions and those in the relevant section. Keep these operating instructions in an accessible place and make them available to the personnel assigned.

### 1.2 Explanation of Symbols

- Triangular warning signs warn of danger.
- Round command symbols tell what to do.

### 1.3 Qualified Personnel

All work such as installation, commissioning, operation, and service of the product may only be carried out by qualified personnel. The personnel must have the necessary qualifications for the corresponding activity and be familiar with the installation, commissioning, operation, and service of the product. The manual and in particular the safety instructions must be carefully read, understood, and observed.

### 1.4 Liability

NTI AG (as manufacturer of LinMot and MagSpring products) excludes all liability for damages and expenses caused by incorrect use of the products. This also applies to false applications, which are caused by NTI AG's own data and notes, for example during sales, support or application activities. It is the responsibility of the user to check the data and information provided by NTI AG for correct applicability in terms of safety. In addition, the entire responsibility for safety-related product functionality lies exclusively with the user. Product warranties are void if products are used with stators, sliders, servo drives, or cables not manufactured by NTI AG unless such use was specifically approved by NTI AG.

NTI AG's warranty is limited to repair or replacement as stated in our standard warranty policy as described in our "terms and conditions" previously supplied to the purchaser of our equipment. Further reference is made to our general terms and conditions.

### 1.5 Copyright

This work is protected by copyright. Under the copyright laws, this publication may not be reproduced or transmitted in any form, electronic or mechanical, including photocopying, recording, microfilm, storing in an information retrieval system, not even for training purposes, or translating, in whole or in part, without the prior written consent of NTI AG.

LinMot® and MagSpring® are registered trademarks of NTI AG.

---

## 2 Safety Instructions

**For your personal safety** — disregarding the following safety measures can lead to severe injury to persons and damage to material:

- Only use the product as directed.
- Never commission the product in the event of visible damage.
- Never commission the product before assembly has been completed.
- Do not carry out any technical changes on the product.
- Only use the accessories approved for the product.
- Only use original spare parts from LinMot.
- Observe all regulations for the prevention of accidents, directives and laws applicable on site.
- Transport, installation, commissioning, and maintenance work must only be carried out by qualified personnel.
- Observe IEC 364 and CENELEC HD 384 or DIN VDE 0100 and IEC report 664 or DIN VDE 0110 and all national regulations for the prevention of accidents.
- According to the basic safety information, qualified, skilled personnel are persons who are familiar with the assembly, installation, commissioning, and operation of the product and who have the qualifications necessary for their occupation.
- Observe all specifications in this documentation.
- This is the condition for safe and trouble-free operation and the achievement of the specified product features.
- The procedural notes and circuit details described in this documentation are only proposals. It is up to the user to check whether they can be transferred to the applications. NTI AG / LinMot does not accept any liability for the suitability of the procedures and circuit proposals described.
- LinMot servo drives, and the accessory components can include live and moving parts (depending on their type of protection) during operation. Surfaces can be hot.
- Non-authorized removal of the required cover, inappropriate use, incorrect installation, or operation create the risk of severe injury to persons or damage to material assets.
- High amounts of energy are produced in the drive. Therefore, it is required to wear personal protective equipment (body protection, headgear, eye protection, hand guard).

**Application as directed:**

- Drives are components, which are designed for installation in electrical systems or machines. They are not to be used as domestic appliances, but only for industrial purposes according to EN 61000-3-2.
- When drives are installed into machines, commissioning (i.e., starting of the operation as directed) is prohibited until it is proven that the machine complies with the regulations of the EC Directive 2006/42/EG (Machinery Directive); EN 60204 must be observed.
- Commissioning (i.e., starting of the operation as directed) is only allowed when there is compliance with the EMC Directive (2014/30/EU).
- The technical data and supply conditions can be obtained from the nameplate and the documentation. They must be strictly observed.

**Transport, storage:**

- Please observe the notes on transport, storage, and appropriate handling.
- Observe the climatic conditions according to the technical data.

**Installation:**

- The drives must be installed and cooled according to the instructions given in the corresponding documentation.
- The ambient air must not exceed degree of pollution 2 according to EN 61800-5-1.
- Ensure proper handling and avoid excessive mechanical stress. Do not bend any components and do not change any insulation distances during transport or handling. Do not touch any electronic components and contacts.
- Drives contain electrostatic sensitive devices, which can easily be damaged by inappropriate handling. Do not damage or destroy any electrical components since this might endanger your health!

**Electrical connection:**

- When working on live drives, observe the applicable national regulations for the prevention of accidents.
- The electrical installation must be carried out according to the appropriate regulations (e.g., cable cross-sections, circuit breakers, fuses, PE connection).
- This product can cause high-frequency interferences in non-industrial environments, which can require measures for interference suppression.

**Operation:**

- If necessary, systems including drives must be equipped with additional monitoring and protection devices according to the valid safety regulations (e.g., law on technical equipment, regulations for the prevention of accidents).
- After the drive has been disconnected from the supply voltage, all live components and power connections must not be touched immediately because capacitors can still be charged. Please observe the corresponding stickers on the drive. All protection covers and doors must be shut during operation.

> ⚠️ **Burn Hazard:** The heat sink (housing) of the drive can have an operating temperature of > 80 °C. Contact with the heat sink results in burns.

> ⚠️ **Risk of Electric Shock:** Before servicing, disconnect supply, wait 5 minutes and measure between PWR+ and PGND to be sure that the capacitors have discharged below 42 VDC. The power terminals Ph1+, Ph1-, Ph2+, Ph2- and PWR+ remain live for at least 5 minutes after disconnecting from the power supplies.

**Grounding:** All metal parts that are exposed to contact during any user operation or servicing and likely to become energized shall be reliably connected to the means for grounding.

---

## 3 System Overview

Typical servo system C1250: Servo drive, motor, and power supply.

*(System overview diagram not reproduced — refer to original PDF for block diagram.)*

---

## 4 Interfaces

*(Interface connector diagram not reproduced — refer to original PDF for connector layout illustration.)*

---

## 5 Functionality

| Feature | C1250-MI-XC-0S | C1250-MI-XC-1S |
|---------|:--------------:|:--------------:|
| Motor Supply: 72 VDC nominal (24...85 VDC) | ● | ● |
| Logic Supply: 24 VDC (22...26 VDC) | ● | ● |
| Motor Phase Current: 25 A peak (0–599 Hz) ¹ | ● | ● |
| Controllable Motors: LinMot P0x- and PR0x- | ● | ● |
| Controllable Motors: Selected motors (contact support) | ● | ● |
| Plug and Play (PnP) Auto Configuration | ● | ● |
| Command Interface: POWERLINK CiA402 | ● | ● |
| Command Interface: PROFINET PROFIdrive | ● | ● |
| Command Interface: Sercos III | ● | ● |
| Command Interface: EtherNet/IP with CIP sync | ● | ● |
| Command Interface: LinUDP | ● | ● |
| Command Interface: EtherCAT CiA402 | ● | ● |
| Command Interface: EtherCAT SoE | not supported | not supported |
| Command Interface: CC-Link | ● | ● |
| Programmable Motion Profiles: up to 100 profiles / 16302 curve points | ● | ● |
| Programmable Command Table: up to 255 entries | ● | ● |
| External Position Sensor: Incremental (RS422, up to 25 Mcounts/s) | ● | ● |
| External Position Sensor: Absolute (SSI, BiSS-B, BiSS-C, EnDat2.1, EnDat2.2) | ● | ● |
| Configuration Interface: RS232 | ● | ● |
| Configuration Interface: Ethernet (EoE, etc.) | ● | ● |
| Integrated Safety Functions: STO (2 Safety Relays) | not supported | ● |
| Calibrated Measuring Functions: Calibrated analog inputs on X4 ² | ●² | ●² |

¹ 28 A peak (0–599 Hz) from firmware release 6.12 and later.  
² Only with the C1250-MI-XC-xS-Cxx type.

---

## 6 Software

The configuration software **LinMot-Talk** is free of charge and can be downloaded from the LinMot homepage.

---

## 7 Power Supply and Grounding

To assure a safe and error free operation, and to avoid severe damage to system components, all system components must be well grounded to protective earth PE. This includes both LinMot and all other control system components on the same ground bus.

Each system component should be tied directly to the ground bus (star pattern). Daisy chaining from component to component is forbidden. (LinMot motors are properly grounded through their power cables when connected to LinMot drives.)

> ⚠️ Power supply connectors must not be connected or disconnected while DC voltage is present. Do not disconnect system components until all LinMot drive LEDs have turned off. (Capacitors in the power supply may not fully discharge for several minutes after input voltage has been disconnected.) Failure to observe these precautions may result in severe damage to electronic components in LinMot motors and/or drives.

> ⚠️ Do not switch Power Supply DC Voltage. All power supply switching and E-Stop breaks should be done to the AC supply voltage of the power supply. Failure to observe these precautions may result in severe damage to the drive.

> **Note:** Inside of the C1250 drive the PWR motor GND and PWR signal GND are connected together and to the GND of the drive housing. It is recommended that the PWR motor GND is **NOT** grounded at another place than inside of the drive to reduce circular currents.

---

## 8 Calibrated Measuring Amplifier (C1250-MI-XC-xS-Cxx)

The drives with the ending `-Cxx` are specially designed for measuring applications. They come with a factory calibration certificate for the analog inputs on X4. The analog inputs on X4 provide a measuring error of less than 1%.

It is the user's responsibility to allow a reasonable period for recalibration. We recommend a calibration interval of **12 months**.

---

## 9 Description of the Connectors / Interfaces

### 9.1 PE — Protective Earth

| Signal | Description |
|--------|-------------|
| PE | Protective Earth |

- Use min. 4 mm² (AWG11)
- Tightening torque: 2 Nm (18 lb·in)

### 9.2 X1 — Motor Supply

*Connector must be ordered separately — see Section 16 Ordering Information.*

| Signal | Description |
|--------|-------------|
| PWR+ | Motor Supply positive |
| PGND | Motor Supply ground |

- Motor Supply: 72 VDC nominal (24...85 VDC)
- Absolute max. rating: 72 VDC +20%
- External Circuit Breaker: **15 A / min. 100 VDC / C-Trip / 5 kA rms SCCR**
- If motor supply voltage exceeds 90 VDC, the drive will go into error state.
- Use 60/75 °C copper conductors only
- Conductor cross-section: 2.5 mm² (AWG14), max. length 3 m

### 9.3 X2/X3 — Motor Connection

**X2 — Motor Phases** *(Connector must be ordered separately)*

| Pin | LinMot Motor | 3-phase EC / Third-party Motor |
|-----|-------------|-------------------------------|
| PH1+ | Motor Phase 1+ (Red) | Motor Phase U (Red) |
| PH1- | Motor Phase 1- (Pink) | Motor Phase V (Pink) |
| PH2+ | Motor Phase 2+ (Blue) | Motor Phase W (Blue) |
| PH2- | Motor Phase 2- (Grey) | Motor Phase X (Grey) |
| PE/SCRN | PE | PE |

- Use 60/75 °C copper conductors only
- Conductor cross-section: 0.5–2.5 mm² (depending on motor current) / AWG 21–14

**X3 — Motor Sensor / Brake** *(DSUB-9 female)*

| Pin | LinMot Motor | EC Motor |
|-----|-------------|---------|
| 1 | Do not connect | DGND |
| 2 | Do not connect | Brake+ |
| 3 | Do not connect | +5 VDC |
| 4 | Do not connect | KTY |
| 5 | +5 VDC | +5 VDC |
| 6 | DGND | DGND |
| 7 | Sensor-Sine | Sensor-Sine / Hall Switch U |
| 8 | Sensor-Cosine | Sensor-Cosine / Hall Switch V |
| 9 | Temp | Hall Switch W |
| Shield | Shield | Shield |

**X3 Notes:**

- Use +5 VDC (X3.5) and DGND (X3.6) only for motor internal hall sensor supply (max. 100 mA).
- Max. motor cable length: 50 m for LinMot Px motors.
- **Brake+:** 24 V / max. 500 mA, peak 1.4 A (will shut down if exceeded); the other terminal must be wired to DGND (X3.1).
- **Caution: Do NOT connect DGND (X3.6) to ground or earth!**
- **Temperature Sensor:** A resistive temperature sensor (PT1000, KTY) can be connected between +5 VDC (X3.2) and KTY (X3.4).
- **Important:** Use Y-style motor cables only (e.g., K15-Y/C). A W-style cable has different shielding and cannot be modified to a Y-style cable.

### 9.4 X4 — Logic Supply / IO Connection

*Spring cage connector — must be ordered separately (Art. No. 0150-3447). See Section 16.*

| Pin | Signal | Description |
|-----|--------|-------------|
| X4.1 | DGND | Logic Ground |
| X4.2 | +24 VDC | Logic Supply 22–26 VDC |
| X4.3 | OUT | Configurable digital Output (can be used as brake output for LinMot motors) |
| X4.4 | OUT | Configurable digital Output |
| X4.5 | IN | Configurable digital Input |
| X4.6 | IN | Configurable digital Input |
| X4.7 | IN | Configurable digital Input |
| X4.8 | IN | Configurable digital Input |
| X4.9 | AnIn | Configurable single-ended analog Input |
| X4.10 | AnIn+ | Configurable differential analog Input positive (with X4.11) |
| X4.11 | AnIn- | Configurable differential analog Input negative (with X4.10) |

**Digital Inputs (X4.5–X4.8):**
- 24 VDC / 5 mA
- Low Level: −0.5 to 5 VDC
- High Level: 15 to 30 VDC

**Digital Outputs (X4.3 & X4.4):**
- 24 VDC / max. 500 mA, peak 1.4 A (will shut down if exceeded)
- Both outputs are high-side switching with integrated pull-down (1.7 kΩ to DGND)

**Analog Inputs:**
- 12-bit A/D converted
- X4.9: 0–10 V, input resistance: > 75 kΩ to DGND
- X4.10/X4.11: ±10 V, input resistance: 28.0 kΩ, common mode range: −5...+10 V to DGND

**Wiring:**
- Use 60/75 °C copper conductors only
- Conductor cross-section max. 1.5 mm²
- Stripping length: 11.5 mm

> ⚠️ **Important:** The 24 VDC logic supply for the control circuit (X4.2) must be protected with an external fuse (**3 A slow blow**).

### 9.5 X13 — External Position Sensor / Differential Hall Switches

*(DSUB-15 female)*

| Pin | ABZ + Hall Switches | SSI / BiSS-B / BiSS-C / EnDat2.1 / EnDat2.2 |
|-----|--------------------|--------------------------------------------|
| 1 | +5 V DC | +5 V DC |
| 2 | A+ | A+ (optional) |
| 3 | A- | A- (optional) |
| 4 | B+ | B+ (optional) |
| 5 | B- | B- (optional) |
| 6 | Z+ | DATA+ |
| 7 | Z- | DATA- |
| 8 | Encoder Alarm (optional) | Encoder Alarm (optional) |
| 9 | DGND | DGND |
| 10 | U+ | nc |
| 11 | U- | nc |
| 12 | V+ | nc |
| 13 | V- | nc |
| 14 | W+ | Clk+ |
| 15 | W- | Clk- |
| case | Shield | Shield |

**Position Encoder Inputs (RS422):**
- Max. counting frequency: 25 Mcounts/s with quadrature decoding
- Minimum 40 ns edge separation must be guaranteed by the encoder under any circumstances
- Max. frequency of each signal: 6.25 MHz

**Differential Hall Switch Inputs (RS422):**
- Input Frequency: < 1 kHz
- Encoder Alarm In: 5 V / 1 mA
- Sensor Supply: 5 VDC, max. 300 mA

### 9.6 X17 – X18 — RealTime Ethernet

| Connector | Function |
|-----------|----------|
| X17 | RT ETH In — RJ-45 |
| X18 | RT ETH Out — RJ-45 |

- 10/100 Mbit/s
- Specification depends on RT Bus type — refer to the corresponding interface documentation.

### 9.7 X19 — System RS232

*(RJ-45)*

| Pin | Signal |
|-----|--------|
| 1 | (Do not connect) |
| 2 | (Do not connect) |
| 3 | RS232 Rx |
| 4 | GND |
| 5 | GND |
| 6 | RS232 Tx |
| 7 | (Do not connect) |
| 8 | (Do not connect) |
| Shield | Shield |

> Use isolated USB-RS232 converter (Art. No. 0150-2473) for configuration over RS232.

### 9.8 X33 — Safety Relays (–1S option only)

*Spring cage connector — must be ordered separately (Art. No. 0150-3451). See Section 16.*

| Pin | Signal | Description |
|-----|--------|-------------|
| 1 / 5 | Ksr f- | Safety Relay 1 / 2 feedback negative |
| 2 / 6 | Ksr f+ | Safety Relay 1 / 2 feedback positive |
| 3 / 7 | Ksr- | Safety Relay 1 / 2 Input negative |
| 4 / 8 | Ksr+ | Safety Relay 1 / 2 Input positive |

**Wiring requirements:**
- Use 60/75 °C copper conductors only
- Conductor cross-section max. 1.5 mm²
- Stripping length: 10 mm
- The state of the feedback contacts **must** be checked after each change of the state of the control contacts.
- Max. current on feedback contacts (Ksr f+ and Ksr f-) must be limited below 1 A.
- **Never connect the safety relays to the logic supply of the drive!**

→ For detailed information see Section 11 Safety Wiring.

### 9.9 S1 – S2 — Address Selectors (DIP Switches)

| Switch | Function |
|--------|----------|
| S1 (bits 5..8) | Bus ID High (0x0…0xF). Bit 5 is LSB, bit 8 is MSB. |
| S2 (bits 1..4) | Bus ID Low (0x0…0xF). Bit 1 is LSB, bit 4 is MSB. |

> Setting both ID High and ID Low to **0xFF** resets the drive to manufacturer default settings.

The use of these switches depends on the type of fieldbus in use — refer to the corresponding interface manual for further information.

### 9.10 System LEDs

| Signal | Color | Description |
|--------|-------|-------------|
| 24VOK | Green | 24 VDC Logic Supply OK |
| EN (enable) | Yellow | Motor Enabled / Error Code Low Nibble |
| WARN | Yellow | Warning / Error Code High Nibble |
| ERROR | Red | Error |

### 9.11 RT Bus LEDs

| Bus Type | L3 (bicolour) | L4 (bicolour) |
|----------|--------------|--------------|
| EtherCAT | RUN (green) | ERR (red) |
| PROFINET | SF (red) | BF (red) |
| POWERLINK | BS (green) | BE (red) |
| EtherNet/IP | MS (green/red) | NS (green/red) |
| SERCOS | S (green/red) | — |
| CC-Link | RUN (green) | ERR (red) |

The blink codes are described in the corresponding interface manuals.

---

## 10 System LED Blink Codes

| ERROR | WARN | EN (enable) | Description |
|-------|------|-------------|-------------|
| OFF | Warning | Operation Enabled | **Normal Operation:** Warnings and operation enabled state are displayed. |
| ON | ● ~2 Hz, 0–15× (Error Code High Nibble) | ● ~2 Hz, 0–15× (Error Code Low Nibble) | **Error:** The error code is shown by a blink code with WARN and EN. The error byte is divided into low and high nibble (4 bits). WARN and EN blink together. The error can be acknowledged. *(Example: WARN blinks 3×, EN blinks 2×; Error Code = 0x32)* |
| ● ~2 Hz | ● ~2 Hz, 0–15× (Error Code High Nibble) | ● ~2 Hz, 0–15× (Error Code Low Nibble) | **Fatal Error:** Same encoding as Error. Fatal errors can only be acknowledged by a reset or power cycle. *(Example: WARN blinks 3×, EN blinks 2×; Error Code = 0x32)* |
| ● ~4 Hz | ● ~2 Hz, 0–15× (Error Code High Nibble) | ● ~2 Hz, 0–15× (Error Code Low Nibble) | **System Error:** Please reinstall firmware or contact support. |
| ● ~0.5 Hz | ● ~0.5 Hz (alternating with ERROR) | Off | **Signal Supply 24 V Too Low:** ERROR and WARN LEDs blink alternating if the signal supply +24 VDC (X4.2) is less than 18 VDC. |
| Off | ○●●● | ●○●● | **Plug & Play Communication Active:** This sequence (WARN on, then EN on, then both off — complete cycle ~1 s) indicates that Plug and Play parameters are being read from the motor. |
| ○● ~4 Hz | ●○ ~4 Hz (alternating) | Off | **Waiting for Defaulting Parameters:** When ID (S1, S2) is set to 0xFF, the drive starts up in a special mode. When the ID is set to 0x00, all parameters will be reset to default. To leave this state, power down the drive and change the ID. See *Usermanual_LinMot-Talk*, chapter Troubleshooting. |
| Off | ○● ~2 Hz | ○● ~2 Hz | **Defaulting Parameters Done:** When parameters have been set to default values (initiated via S1/S2 on power-up), the WARN and EN LEDs blink together at 2 Hz. To leave this state, power down the drive. See *Usermanual_LinMot-Talk*, chapter Troubleshooting. |

The meaning of error codes can be found in **Usermanual_MotionCtrl_Software_SG5-SG7** and the user manual of the installed interface software. These documents are provided with LinMot-Talk and can be downloaded from [www.linmot.com](https://www.linmot.com).

---

## 11 Safety Wiring

The C1250 drives with the **–1S option** have internal safety functions: Two safety relays Ksr in series, which support the supply voltage for the motor drivers (normally open). There is also a feedback contact for each relay (normally closed).

To enable the –1S drives, **both relays must be switched on**.

**Minimal wiring:**
- Connect X33.8 and X33.4 to 24 VDC (from safety supply)
- Connect X33.7 and X33.3 to GND (from safety supply)
- **Never connect X33.8 and X33.4 to the logic supply of X4!**
- **Never disconnect X33 when the STO is powered!**

> The overvoltage protection must be provided externally and sized according to the safety circuit of the machine. The voltage on any pin of X33 must be limited below 100 V referenced to DGND.

> The drop-out time of the relays depends on the external circuitry.

**Safety Relay Ksr Specifications:**

| Parameter | Value |
|-----------|-------|
| Nominal voltage | 24 VDC |
| Min. pick-up voltage at 20 °C | ≤ 16.8 V |
| Drop-out voltage at 20 °C | ≥ 2.4 V |
| Drop-out time (no protection circuit) | Typically 3 ms |
| Coil resistance at 20 °C | 2 100 Ω ± 10% |
| Type | EN 50205, type A — Relay with forcibly guided contacts per IEC 61810-3 |
| Contact lifetime | > 10 000 000 operations |
| Manufacturer and type | Elesta relays / SIS112 24VDC |
| Max. current on feedback contacts (Ksr f+ and Ksr f-) | < 1 A |

**Drive Classification per EN ISO 13849-1 (Safety of Machinery):**

| Parameter | Value |
|-----------|-------|
| Category | cat. 3 |
| Performance Level | PL = d |
| Diagnostic Coverage | DC = high (99%) |
| Mean Time to Hazardous Failure (MTTFd) | high (100 years typically — see calculation below) |

DC (Diagnostic Coverage) is high (99%) assuming that the state of the feedback contacts is checked after each change of the state of the control contacts.

**Example MTTFd Calculation:**

Assuming the safety function is requested every 20 s on a machine running 24 h/day, 7 days/week:

| Parameter | Value |
|-----------|-------|
| B10 | 10 000 000 |
| B10d | 20 000 000 (per EN ISO 13849-1:2008 table C.1) |
| nop | (24 h/day × 365.25 days/year × 3600 s/h) / 20 s = **1 577 880 operations/year** |
| MTTFd | B10d / (0.1 × nop) = **126.75 years** (limited to 100 years per standard) = high |

---

## 12 Physical Dimensions

**C1250-MI Series — Single Axis Drive**

| Parameter | Unit | C1250-MI-XC-0S | C1250-MI-XC-1S |
|-----------|------|:--------------:|:--------------:|
| Width | mm (in) | 25.3 (1.0) | 25.3 (1.0) |
| Height | mm (in) | 166 (6.54) | 176 (6.93) |
| Height with fixings | mm (in) | 206 (8.11) | 216 (8.50) |
| Depth | mm (in) | 106 (4.2) | 106 (4.2) |
| Weight | g (lb) | 630 (1.4) | 700 (1.54) |
| Mounting Screws | | 2 × M5 | 2 × M5 |
| Mounting Distance | mm (in) | 198 (7.80) | 198 (7.80) |
| Degree of Protection | IP | 20 | 20 |
| Storage Temperature | °C | −25…40 | −25…40 |
| Transport Temperature | °C | −25…70 | −25…70 |
| Operating Temperature | °C | 0…40 at rated data (40…50 with power derating*) | 0…40 at rated data (40…50 with power derating*) |
| Relative Humidity | | < 95% (non-condensing) | < 95% (non-condensing) |
| Pollution | | IEC/EN 60664-1, Pollution degree 2 | IEC/EN 60664-1, Pollution degree 2 |
| Shock resistance (16 ms) | | 3.5 g | 3.5 g |
| Vibration resistance (10–200 Hz) | | 1 g | 1 g |
| Max. Case Temperature | °C | 70 | 70 |
| Max. Power Dissipation | W | 30 | 30 |
| Mounting place | | In the control cabinet (at least IP54) | In the control cabinet (at least IP54) |
| Mounting position | | Vertical | Vertical |
| Distance between Drives (without derating*) | mm (in) | 20 (0.8) horizontal / 50 (2.0) vertical | 20 (0.8) horizontal / 50 (2.0) vertical |
| Distance between Drives (with derating*) | mm (in) | 5 (0.2) horizontal / 20 (0.8) vertical | 5 (0.2) horizontal / 20 (0.8) vertical |

*\* The derating depends on the cabinet situation. The temperature of the drive should be checked under full load (the temperature should be stable, which may take an hour or more). This verifies that enough margin exists when the cabinet reaches its maximum allowable temperature of 40 °C. The warning level of the drive is configured by default to 75 °C and the error level to 80 °C.*

---

## 13 Power Supply Requirements

### 13.1 Motor Power Supply

The calculation of the needed power for the motor supply depends on the application and the motor used. The nominal supply voltage is **72 VDC**. The possible range is from **24 to 85 VDC**.

> ⚠️ The motor supply can rise to 95 VDC when braking. This means that everything connected to that power supply needs a dielectric withstand voltage of at least 100 VDC (additional capacitors, etc.). Due to high braking voltage and sudden load variations of linear motor applications, only compatible power supplies can be used (see Section 16 Ordering Information).

### 13.2 Signal Power Supply

The logic supply needs a regulated power supply with a nominal voltage of 24 VDC. The voltage must be between 22 and 26 VDC.

**Current drawn from logic supply:**

| Condition | Current |
|-----------|---------|
| Minimum (no load on outputs) | 0.5 A |
| Typical (all 2 outputs "on" with 100 mA load, brake with no load) | 0.7 A |
| Maximum (all 2 outputs "on" with 500 mA peak load, brake with 500 mA peak load) | 2.0 A |

> ⚠️ Do not connect the safety relays to the 24 VDC Signal Supply! Use a separate power supply for the safety circuit.

> ⚠️ The 24 VDC supply for the control circuit must be protected with an external fuse (**3 A slow blow**).

---

## 14 Regeneration

If the power supply rises too high during braking, connect an additional capacitor to the motor power supply. It is recommended to use a capacitor **≥ 10 000 μF** (install capacitor close to the drive supply).

---

## 15 Safety Notes for Installation According to UL

**Markings:**
- **Wiring terminal markings:** See markings on the enclosure and the corresponding chapters in the installation guide.
- **Cautionary Marking:** See markings on the enclosure and the corresponding chapters in the installation guide.
- The transients must be limited to max. **0.8 kV** on the line side of the drive.
- The 24 VDC supply for the control circuit must be protected with an external **UL Listed 3 A DC fuse**.
- A separate 24 VDC power supply protected with an external **UL Listed 3 A fuse** connected to the output of the power supply must be used to protect the secondary control circuit (safety relays on X33).
- Integral solid-state short-circuit protection does not provide branch circuit protection. Branch circuit protection must be provided in accordance with the National Electrical Code and any additional local codes.

**Ratings for cULus Listing:**

| Parameter | Value |
|-----------|-------|
| Input Voltage | 72 VDC |
| Input Current | 6.5 A |
| Output Voltage | 43 V rms |
| Output Current | 5 A rms |
| Number of Phases | 2–3 |
| Frequency Range | 0–599 Hz |
| Duty Cycle Rating | 4–96% |
| Relays — Rated Contacts (–1S variant only) | max. 24 VDC, 6 A |
| Relays — Coil | 24 VDC |
| Relays — Surrounding Air Temperature | max. 85 °C |
| Control Power (X4.2) | 24 VDC (protected with external UL Listed 3 A fuse) |
| Control Power — Surrounding Air Temperature | max. 50 °C |

Suitable for use on a circuit capable of delivering not more than **5 kA DC, 72 VDC maximum**.  
(Branch circuit protection on X1: External Circuit Breaker: 15 A / min. 100 VDC / C-Trip / 5 kA rms SCCR)

---

## 16 Ordering Information

### 16.1 Drives

**Standard Drives:**

| Part Number | Description | Art. No. |
|-------------|-------------|----------|
| C1250-MI-XC-0S-000 | Multi Interface Drive (72V/25A) | 0150-5591 |
| C1250-MI-XC-1S-000 | Multi Interface Drive (72V/25A), STO | 0150-5589 |
| C1250-MI-XC-1S-0PD | Multi Interface Drive (72V/25A), STO, PROFINET/PROFIdrive installed | 0150-5732 |
| C1250-MI-XC-1S-0CM | Multi Interface Drive (72V/25A), STO, EtherNet/IP CIP Sync installed | 0150-5733 |
| C1250-MI-XC-1S-0LU | Multi Interface Drive (72V/25A), STO, LinUDP installed | 0150-5734 |
| C1250-MI-XC-1S-0PL | Multi Interface Drive (72V/25A), STO, POWERLINK installed | 0150-5735 |
| C1250-MI-XC-1S-0SC | Multi Interface Drive (72V/25A), STO, Sercos III installed | 0150-5736 |
| C1250-MI-XC-1S-0DS | Multi Interface Drive (72V/25A), STO, EtherCAT/CiA402 installed | 0150-5737 |
| C1250-MI-XC-1S-0CC | Multi Interface Drive (72V/25A), STO, CC-Link installed | 0150-5738 |
| C1250-MI-XC-0S-0PD | Multi Interface Drive (72V/25A), PROFINET/PROFIdrive installed | 0150-5746 |
| C1250-MI-XC-0S-0CM | Multi Interface Drive (72V/25A), EtherNet/IP CIP Sync installed | 0150-5747 |
| C1250-MI-XC-0S-0LU | Multi Interface Drive (72V/25A), LinUDP installed | 0150-5748 |
| C1250-MI-XC-0S-0PL | Multi Interface Drive (72V/25A), POWERLINK installed | 0150-5749 |
| C1250-MI-XC-0S-0SC | Multi Interface Drive (72V/25A), Sercos III installed | 0150-5750 |
| C1250-MI-XC-0S-0DS | Multi Interface Drive (72V/25A), EtherCAT/CiA402 installed | 0150-5751 |
| C1250-MI-XC-0S-0CC | Multi Interface Drive (72V/25A), CC-Link installed | 0150-5752 |

**Calibrated Drives (factory-calibrated analog inputs on X4, <1% error):**

| Part Number | Description | Art. No. |
|-------------|-------------|----------|
| C1250-MI-XC-0S-C00 | Multi Interface Drive (72V/25A), Calibrated Measuring Amplifier | 0150-5592 |
| C1250-MI-XC-1S-C00 | Multi Interface Drive (72V/25A), STO, Calibrated Measuring Amplifier | 0150-5590 |
| C1250-MI-XC-1S-CPD | Multi Interface Drive (72V/25A), STO, Calibrated, PROFINET/PROFIdrive installed | 0150-5725 |
| C1250-MI-XC-1S-CCM | Multi Interface Drive (72V/25A), STO, Calibrated, EtherNet/IP CIP Sync installed | 0150-5726 |
| C1250-MI-XC-1S-CLU | Multi Interface Drive (72V/25A), STO, Calibrated, LinUDP installed | 0150-5727 |
| C1250-MI-XC-1S-CPL | Multi Interface Drive (72V/25A), STO, Calibrated, POWERLINK installed | 0150-5728 |
| C1250-MI-XC-1S-CSC | Multi Interface Drive (72V/25A), STO, Calibrated, Sercos III installed | 0150-5729 |
| C1250-MI-XC-1S-CDS | Multi Interface Drive (72V/25A), STO, Calibrated, EtherCAT/CiA402 installed | 0150-5730 |
| C1250-MI-XC-1S-CCC | Multi Interface Drive (72V/25A), STO, Calibrated, CC-Link installed | 0150-5731 |
| C1250-MI-XC-0S-CPD | Multi Interface Drive (72V/25A), Calibrated, PROFINET/PROFIdrive installed | 0150-5739 |
| C1250-MI-XC-0S-CCM | Multi Interface Drive (72V/25A), Calibrated, EtherNet/IP CIP Sync installed | 0150-5740 |
| C1250-MI-XC-0S-CLU | Multi Interface Drive (72V/25A), Calibrated, LinUDP installed | 0150-5741 |
| C1250-MI-XC-0S-CPL | Multi Interface Drive (72V/25A), Calibrated, POWERLINK installed | 0150-5742 |
| C1250-MI-XC-0S-CSC | Multi Interface Drive (72V/25A), Calibrated, Sercos III installed | 0150-5743 |
| C1250-MI-XC-0S-CDS | Multi Interface Drive (72V/25A), Calibrated, EtherCAT/CiA402 installed | 0150-5744 |
| C1250-MI-XC-0S-CCC | Multi Interface Drive (72V/25A), Calibrated, CC-Link installed | 0150-5745 |

> There are drives with a certain interface preinstalled. On -MI drives, any available interface can nevertheless be installed during firmware installation.

### 16.2 Accessories

| Part Number | Description | Art. No. |
|-------------|-------------|----------|
| **DC01-C1X00-0S/X1/X4** | **Drive Connector Set for C1X00-0S** | **0150-3527** |
| **DC01-C1X00-1S/X1/X4/X33** | **Drive Connector Set for C1X00-1S** | **0150-3528** |
| DC01-C1X00/X1 | Drive Connector for PWR 72 VDC Input | 0150-3525 |
| DC01-C1X00/X2 | Drive Connector Motor Phases | 0150-3526 |
| **DC01-Signal/X4** | **Drive Connector 24 VDC & Logic** | **0150-3447** |
| DC01-Safety/X33 | Drive Connector Safety | 0150-3451 |
| **Isolated USB-RS232 converter** | **Isolated USB RS232 converter with config. cable** | **0150-2473** |
| Isolated USB-serial converter | Isolated USB RS232/422/485 converter | 0150-3120 |
| Recalibration Service | Calibration Drive Series C1200 (analog inputs on X4 of C1250-xx-XC-xS-Cxx) | 0150-4164 |

**Compatible Power Supplies:**

| Part Number | Description | Art. No. |
|-------------|-------------|----------|
| **S02-72/1000** | **Power Supply 72 V/1000 W, 3×400–480 VAC** | **0150-4535** |
| S01-72/1000 | Power Supply 72 V/1000 W, 3×340–550 VAC | 0150-1872 |
| S01-72/500 | Power Supply 72 V/500 W, 1×120/230 VAC | 0150-1874 |
| S01-24/500 | Power Supply 24 V/500 W, 1×120/230 VAC | 0150-2480 |
| S01-48/300 | Power Supply 48 V/300 W, 1×120/230 VAC | 0150-1941 |
| S01-48/600 | Power Supply 48 V/600 W, 1×120/230 VAC | 0150-1946 |
| T01-72/420-Multi | T-Supply 72 V/420 VA, 3×230/400/480 VAC | 0150-1869 |
| T01-72/900-Multi | T-Supply 900 VA, 3×230/400/480 VAC | 0150-1870 |
| T01-72/1500-Multi | T-Supply 1500 VA, 3×230/400/480 VAC | 0150-1871 |
| T01-72/420-1ph | T-Supply 420 VA, 1×208/220/230/240 VAC | 0150-1859 |

> **Bold items are strongly recommended accessories.**  
> The connectors must be ordered separately and are **not included** with the drive.  
> Use 0150-2473 (isolated USB RS232 converter) for configuration.

---

## 17 International Certifications

| Region | Certification |
|--------|--------------|
| Europe | CE Marking — see Section 17.3 |
| UK | UKCA Marking — see Section 17.4 |
| IECEE CB Scheme | Ref. Certif. No. CH-11687 |
| USA / Canada | UL Listed — File number E316095; UL 508C Power Conversion Equipment; CSA C22.2 Industrial Control Equipment — see Section 17.2 |

All products marked with the UL symbol are tested and listed by Underwriters Laboratories; production facilities are checked quarterly by a UL inspector. This mark is valid for the USA and Canada.

### 17.1 IECEE CB Scheme — CB Test Certificate

*(Certificate document — refer to original PDF for the full certificate pages.)*

**Summary of compliance with EN standards (CENELEC countries):**

- EN IEC 61800-3:2018
- EN 55011:2016+A1:2017+A11:2020+A2:2021 Class A
- EN 55032:2015+A1:2020+A11:2020 Class A
- EN 61000-3-2:2014
- EN IEC 61000-3-2:2019+A1:2021
- EN 61000-3-3:2013+A1:2019+A2:2021

**Summary of compliance with other IEC standards:**

- CISPR 11:2015 Class A (and AMD1:2016, AMD2:2019)
- CISPR 32:2015 Class A (and AMD1:2019)
- IEC 61000-3-2:2018 (ed.5) and AMD1:2020
- IEC 61000-3-3:2013 (ed.3), AMD1:2017, and AMD2:2021

### 17.2 UL Listing

*(UL listing certificate — refer to original PDF for the full certificate page.)*

### 17.3 EU Declaration of Conformity — CE Marking

NTI AG / LinMot®  
Bodenaeckerstrasse 2, 8957 Spreitenbach, Switzerland  
Tel.: +41 (0)56 419 91 91 | Fax: +41 (0)56 419 91 92

Declares under sole responsibility the compliance of the products:

- Drives of the Series **C1250-MI-XC-xS-xxx**

with the **EMC Directive 2014/30/EU**.

Applied harmonized standards:
- EN 61800-3:2004 + A1:2012
- EN 61800-3:2018

According to the EMC directive, the listed devices are not independently operable products. Compliance of the directive requires the correct installation of the product, the observance of specific installation guides and product documentation. This was tested on specific system configurations.

The product must be mounted and used in strict accordance with the installation instructions contained within the installation guide, a copy of which may be obtained from NTI AG.

**Company:** NTI AG  
**Date:** Spreitenbach, 04.08.2022  
**Signed:** Dr. Ronald Rohner / CEO NTI AG

### 17.4 UK Declaration of Conformity — UKCA Marking

NTI AG / LinMot®  
Bodenaeckerstrasse 2, 8957 Spreitenbach, Switzerland  
Tel.: +41 (0)56 419 91 91 | Fax: +41 (0)56 419 91 92

Declares under sole responsibility the compliance of the products:

- Drives of the Series **C1250-MI-XC-xS-xxx**

with the **EMC Regulation S.I. 2016 No. 1091**.

Applied designated standards:
- EN 61800-3:2004 + A1:2012
- EN 61800-3:2018

According to the EMC regulation, the listed devices are not independently operable products. Compliance requires the correct installation of the product, the observance of specific installation guides and product documentation. This was tested on specific system configurations.

The product must be mounted and used in strict accordance with the installation instructions contained within the installation guide.

**Company:** NTI AG  
**Date:** Spreitenbach, 23.03.2022  
**Signed:** Dr. Ronald Rohner / CEO NTI AG

---

## Contact Information

**Europe / Asia Headquarters:**  
NTI AG — LinMot & MagSpring  
Bodenaeckerstrasse 2, CH-8957 Spreitenbach, Switzerland  
Sales/Administration: +41 56 419 91 91 | office@linmot.com  
Technical Support: +41 56 544 71 00 | support@linmot.com  
Web: https://www.linmot.com

**North / South America Headquarters:**  
LinMot USA Inc.  
N1922 State Road 120, Unit 1, Lake Geneva, WI 53147, USA  
Sales/Administration: 262.743.2555 | usasales@linmot.com  
Technical Support: 262.743.2555 | usasupport@linmot.com  
Web: https://www.linmot-usa.com

Find a distributor: https://linmot.com/contact/
