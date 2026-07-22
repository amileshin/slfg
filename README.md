# slfg

[![CI](https://github.com/amileshin/slfg/actions/workflows/ci.yml/badge.svg)](https://github.com/amileshin/slfg/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/amileshin/slfg/branch/main/graph/badge.svg)](https://codecov.io/gh/amileshin/slfg)
[![Coverage Status](https://coveralls.io/repos/github/amileshin/slfg/badge.svg?branch=main)](https://coveralls.io/github/amileshin/slfg?branch=main)

Библиотека логирования для Go, переносящая модель SLF4J + Logback в идиоматичный Go.

Идея — дать Go-проектам, мигрирующим с Java, привычную систему логирования: иерархия логгеров по точкам, конфигурация из `logback.xml` / `logback.yaml`, MDC, маркеры, appenders с фильтрами и pattern layout — без необходимости переосмыслять логирование с нуля.
