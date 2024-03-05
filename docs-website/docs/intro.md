---
sidebar_position: 1
slug: /
title: Introduction
---

# SpecPipe: A scalable AI/ML-facilitating data pipeline for spectrum

By: Omair Alam, Will Almy, Alice Lee and Huo-Ming Shu

---

In today's interconnected world, radio spectrum signals surround us, yet there exist noticeable limitations in the data systems created to access, monitor, perform AI experiments, and contribute to this analog data.

To democratize the access and usage of spectrum data, we have built [SpecPipe](https://github.com/ml4wireless/specpipe), a distributed AI/ML data pipeline. This platform’s core values of accessibility, extensibility and scalability ensure that individual users can start to work with radio data with inexpensive hardware, minimal configuration, and a smooth onboarding process.

We have accomplished this goal of improving access to spectrum data by building SpecPipe as an open-source project free for people to access and use, with easy to follow documentation, and a plethora of startup examples that allow users to understand our framework interactively.

For more details on the architecture of SpecPipe [click here](./architecture).

---

## Setup Guide

### Setting up SpecPipe Server and Edge Device

Note: This only applies if you want to setup to send Radio Data as an [Edge Node](./architeture/#edge-nodes). If you only want to build an application that accesses data, you don't need to have an SDR.

#### 1. Setting up NATS on Amazon Web Services

Note: This is only required if you want to deploy your server on AWS. If you want a local NATS installation, [step 3](./#3-setting-up-specpipe-edge-device) will guide you through that.

The [NATS AWS tutorial](./NATS.md) tutorial will help you setup NATS on AWS that you will be able to then use when setting up your SpecPipe Edge devices in [step 3](./#3-setting-up-specpipe-edge-device).

#### 2. Connecting a Software Defined Radio (SDR) to a machine

The SDR will be collecting data in real-time.

For the hardware setup, see: https://www.youtube.com/watch?v=uM8NkB2nIis

You can follow this [software setup](https://www.youtube.com/watch?v=bT2WZhKBkRk) to confirm that your SDR is functioning correctly. SpecPipe doesn't make use of this software so you won't need it for anything besides this.

#### 3. Setting up SpecPipe Edge Device

Please follow [these instructions](https://github.com/ml4wireless/specpipe?tab=readme-ov-file#getting-started) to setup an Edge device for sending Radio Data using the SDR added in Step 1.

#### 4. Creating Monitoring Dashboard

To create a monitoring dashboard for your system, follow [this](./dashboard.md) tutorial.

### Building Specpipe Applicataion

This step only applies if you would like to build a [Specpipe Application](./architecture#applications)

Please follow [these instructions](https://github.com/ml4wireless/specpipe-sdk-py?tab=readme-ov-file#installation--usage) to create a SpecPipe Application.



