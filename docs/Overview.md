## uCISC Programming Guide

1. [Getting Started](programming_guide/1.0_Getting_Started.md)
   1. [Configuring IntelliJ](programming_guide/1.1_Configuring_IntelliJ.md)
   2. [Configuring VIM](programming_guide/1.2_Configuring_VIM.md)
   3. [Compiling uCISC Code](programming_guide/1.3_Compiling_uCISC_Code.md)
   4. [Using the uCISC Simulator](programming_guide/1.4_Simulating_uCISC.md)
   5. [Using the uCISC Soft Core](programming_guide/1.5_Running_uCISC_Soft_Core.md)
2. [Introduction to Programming With uCISC](programming_guide/2.0_Program_With_uCISC.md)
   1. [Accessing External Devices](programming_guide/2.1_Accessing_Devices.md)
   2. [Common Devices](programming_guide/2.2.0_Common_Devices.md)
      1. [GPIO](programming_guide/2.2.1_GPIO_Devices.md)
      2. [I2C](programming_guide/2.2.2_I2C_Devices.md)
      3. [UART](programming_guide/2.2.3_UART_Devices.md)
      4. [Video Devices](programming_guide/2.2.4_Video_Devices.md)
   3. [Advanced uCISC Programming Techiques](programming_guide/2.3_Advanced_Programming_Techniques.md)
3. [uCISC Syntax Quick Reference](programming_guide/3_Syntax_Quick_Reference.md)
4. [Standard Libraries]()
5. [Instruction Set Details](programming_guide/5_Instruction_Set_Details.md)

The [GrokThis YouTube channel](https://www.youtube.com/c/GrokThis) may also be helpful.

### uCISC is for People

> What would a computer optimized for people be like?

A machine optimized for processing still has to be programmed by people. How
have we bridged that gap? With tools. We have built a complex set of tools
(compilers, languages, IDEs and more) that take our human words and translate
them into machine processable bits.

uCISC (pronounced like micro-sisk) is my attempt at exploring what happens
when you optimize for people first, even at the processor level. You can't
get away from 1's and 0's since they are fundamental to transistor logic
based machines, but let's organize those bits into something that works
better for us as human beings.

## How to put People First

To me, a people first computer means:

1. It has a simple mental model
2. It is fully featured
3. It has minimal translation layers
4. It is usable at human scale

### Simple Mental Model

A single person can understand the entire thing relatively easily.

*People should think:*
  * "Oh, I get it!"
  * "I could do this!"
  * "This is too simple, are you sure this is fully featured?"
  * "Obviously, that would be ..." (when asked how a new hardware feature would work)

*People should NOT think:*
  * "That's too complicated"
  * "I'm having a hard time understanding it all."

### Fully Featured

With enough development resources, it could be used to build a fully functional
computer as a daily driver.

*People should think:*
  * "Wow, I see how you really could do anything with this!"
  * "I wish we had more programmers building software for this"
  * "Let's scale this up and see where it takes us!"

*People should NOT think:*
  * "You can't do X with this" (assuming they are right)

### Minimal translation layers

A programmer should feel connected to the hardware rather than it being
abstracted away. The computer should feel very transparent, almost as if
it is asking you to look under the hood.

*People should think:*
  * "With a bit of work, I could write down the exact instructions that this is using"
  * "I don't know how this works, but I can just open it up and look"

*People should NOT think:*
  * "I have no idea how this actually runs on the hardware"
  * "I just treat this as a black box that does magic for me"

### Usable at human scale

*People should think:*
  * "I'm going to build a library for uCISC that does X for other people to use"
  * "Let's build uCISC software systems"
  * "If we agree on X, we can work independently put it together at the end"

*People should NOT think:*
  * "I need C/python/etc in order to do anything more than hello world"
  * "If I try to build this on a uCISC machine, I'll never finish"

You may choose different criteria for your human first machine.
If you are interested in why I chose these criteria, you can start at the
[Beginning](history/1_Beginnings.md) of my personal story. You can also
follow the ongoing saga of actually building this on the
[GrokThis YouTube Channel](https://www.youtube.com/channel/UCh4OpfF7T7UtezGejRTLxCw).

### Special Thanks

I must give vast amounts of credit to Kartik Agaram. I started down this journey
independently, but he has been instrumental in helping me come so far so quicky.
At some point a colleague of my saw one of his blog posts on hacker news about
[Mu: A minimal hobbyist computing stack](http://akkartik.name/post/mu-2019-1)
and recognized I would be interested. The rest, as they say, is history, and I
am unbelievably grateful to his advice, critiques and motivating influence.

## uCISC History and Motivation

1. [Beginnings](1_Beginnings.md)
2. [Transparent Computing](2_Transparent_Computing.md)
3. [Rough Hardware Outline](3_Rough_Hardware_Outline.md)