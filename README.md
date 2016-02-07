# gous-vide
## What? 
An all-in-one application to run my (your?) Raspberry Pi-based Sous Vide cooker. A ~~[PID controller](https://en.wikipedia.org/wiki/PID_controller)~~ **simple polling routine** to maintain temperature, REST services to control state, target temperature, sensors, etc. and a pretty front-end to consume said services. From a hardware standpoint, I was able to find some outstanding resources online to get to a point where I can read temp from a [DS18B20 "1-Wire" Temperature Sensor](https://learn.adafruit.com/adafruits-raspberry-pi-lesson-11-ds18b20-temperature-sensing/hardware) and [control an RF power outlet](http://timleland.com/wireless-power-outlets/). As a bonus, I had an excuse to buy and play with a soldering iron, so I've already seen a considerable return on the time commitment thus far.  I know that's not much to go off of - I'll put together an .io page with some more details on how to turn your $35 Pi into a mean water bath cooking machine once I've got the software in a good state. 

## Why?
Yeah, I know [Anova](http://anovaculinary.com/anova-precision-cooker/) already makes a great Sous Vide. It has a nice mobile app, it has bluetooth/wifi, and it's pretty. But cooking delicious meats is only a secondary goal of this project. I have *at least* a few other neglected items on my dusty TODO list of personal projects, coincidentally sitting right next to that line item to use the Raspberry Pi that I got for Christmas last year:
* Lots of neglected TODO's
* ~~**Learn Clojure**~~ See [ClojousVide](https://github.com/Mattmanx/clojous-vide) for a rundown of that project, and why I switched to Golang
* **Learn Golang**
* **Learn ReactJS**
* **Use Raspberry Pi**
* Lots of neglected TODO's

So, I thought, why not kill 2 (3?) birds with one stone and incorporate several into a single project? 

## How? 
This is a work in progress, but here's the plan: 
* Go Components
 * ~~**PID Controller**~~ - I initially created a time-centric PID controller in Clojure using Stian Eikeland's excellent [blog post](http://blog.eikeland.se/2014/10/06/pid-transducer/) for inspiration.  But I soon realized that, given my hardware setup, binary heater control, and the high heat capacity of water, I could create a simple ticked poll mechanism that compared the current temperature to target and turned the heaters on / off accordingly.  Not nearly as fun as creating a PID controller, but simple and has worked just fine for me.     
 * **Go Routines** to poll temperature, run a _program_ to maintain a target temperature 
 * **Embedded BoltDB Datasource** to persist temperature history, heater history, programs and recipes.   
 * **REST Services** to allow complete control and visibility over HTTP, using Gorilla's [Mux](https://github.com/gorilla/mux) router for convenience
* Front-End Components
 * A **ReactJS/Redux front-end** to consume the REST services and provide a UI for controlling and tracking the Sous Vide.
 * A datavis library integrated with ReactJS to chart temperature history, heater history, etc. with pretty visuals

# Getting Started
## Hardware & Material Prereqs
TBD
## Building the Sous Vide
TBD
## Running Gous-Vide
TBD
