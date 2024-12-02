# Bluetooth Speed and Cadence Bike Sensor Client
Use a Bluetooth Speed and Cadence Sensor and display stats such as current speed and pedal cadence as well as total distance and duration of the exercise.

The idea of the project is to be able to display workout statistics on a screen / HUD instead of having to monitor a tiny bike-computer screen.

## Estimating Power output
Trainer used: [Elite Novo Mag Force 8](https://www.amazon.de/Elite-Rollentrainer-Novo-Mag-Force/dp/B01K52T51M)
#### Setting: using level 3 of 8
![Elite Novo Mag Force, level 3/8](./TrainerPower.png)

In order to be able to provide the user with an estimate of their power output, we used the values obtained from the official [Elite App](https://www.elite-it.com/en/products/app-software/my-e-training) and integrated them. We used the level 3rd of 8 for the tracking, but this can be adjusted in future versions (e.g by passing the level setting from the WebApp to the backend).
The power output estimation is needed to be able to allow the users to track their training ride properly (needed for total energy calculation as well as calories burnt) and will be crucial when exporting the training to [Strava](https://www.strava.com) or similar platforms.

## Resources
* [Speed and Cadence Data Definition](https://github.com/sputnikdev/bluetooth-gatt-parser/blob/master/src/main/resources/gatt/characteristic/org.bluetooth.characteristic.csc_measurement.xml)
* [Gauge JS Component](https://bernii.github.io/gauge.js/#!)
* [Websockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API)

## Possible Extensions
* allow exporting a training session
* allow changing the difficulty level on the trainer