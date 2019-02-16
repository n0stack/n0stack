.. n0stack documentation master file, created by
   sphinx-quickstart on Sun Feb 17 02:00:07 2019.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.


n0stack documentation
===================================

The n0stack is a simple cloud provider using gRPC.

Description
===========

The n0stack is...

- a cloud provider.
    - You can use some features: booting VMs, managing networks and so on (see also /n0proto.)
- simple.
    - There are shortcode and fewer options.
- using gRPC.
    - A unified interface increase reusability.
- able to be used as library and framework.
    - You can concentrate to develop your logic by sharing libraries and frameworks for middleware, test, and deployment.

Motivation
==========

Cloud providers have various forms depending on users.
This problem has been solved with many options and add-ons (e.g. OpenStack configuration file is very long.)
However, it is difficult to adapt to the application by options, then it is necessary to read or rewrite long abstracted codes.
Therefore, I thought that it would be better to code on your hands from beginning.

There are some problems to develop cloud providers from scratch: no library, software quality, man-hour, and deployment.
The n0stack wants to solve such problems.

.. ## Demo

.. TODO: READMEからインポート、現時点では相対パスやバッジなどが面倒
.. .. mdinclude:: ../README.md

.. toctree::
    :caption: User Documentation
    :maxdepth: 1
    :glob:

    user/quick_start.md
    user/*
    user/usecases/README.rst

.. toctree::
    :caption: Developer Documentation
    :maxdepth: 1
    :glob:

    developer/adr/README.rst
    developer/*
