# teltonika2tak-cot

# Teltonika GPS Tracker to TAK Server

This repository contains a server application built in Golang to receive and process data from Teltonika GPS trackers (for example, model TMT250), and send it to a TAK server in real time. The application uses the [Teltonika parser](https://github.com/filipkroca/teltonikaparser) for processing GPS tracker data and converts it to a format compatible with the Cursor on Target (CoT) standard for display on a TAK server.

## Features

1. Golang server for receiving data from Teltonika GPS trackers.
2. Conversion of received data into CoT packets for display on a TAK server.
3. Real-time data transmission to a TAK server.
4. Configuration files that enable the system to run in a Docker container and specify the address of the WINTAK server.
5. A GUI built with Electron that displays the time, coordinates, direction, speed, and other parameters of the data sent by the trackers. It also allows configuring a conversion table from the tracker's IMEI to the Cursor on Target.

## API Usage

### Teltonika Parser

The Teltonika parser is used to parse and validate data structures from Teltonika UDP packets. The API includes types such as `Decoded` and `AvlData`, and functions such as `Decode()`.

For more information, please refer to the [official repository](https://github.com/filipkroca/teltonikaparser).

## Installation and Configuration

Details about installation, configuration, and usage will be provided soon.

## Note

Detailed API documentation for Cursor on Target (CoT) is still being gathered. Please refer to the official CoT project documentation or community resources for more information.

## Contributing

Contributions are welcome! Please read our Contributing guidelines for details.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
