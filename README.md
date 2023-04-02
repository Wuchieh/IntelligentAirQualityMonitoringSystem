# Intelligent Air Quality Monitoring System v2
English | [繁體中文](README_tw.md)

<a href="https://github.com/wuchieh"><img src="imgs/LogoWhile.png" alt="wuchieh" style="height: 100px;"></a>
&emsp;
<a href="https://github.com/MeowXiaoXiang"><img src="https://github.com/MeowXiaoXiang.png" alt="MeowXiaoXiang" style="height: 100px;"></a>

## Require
```
PostgreSQL
Redis
Nginx
go version > 1.20
```

## Introduction
### This is the second generation of `Intelligent Air Quality Monitoring System`.
### Most of the system has been optimized
### Introduced front-end and back-end separation in the new version
### The front-end uses the Vue3 framework instead of rendering on the back-end to reduce back-end stress
### Also changed the database from MySQL to [PostgreSQL](https://github.com/lib/pq)
### And introduced [Redis](https://github.com/redis/go-redis/) to improve some efficiency and adapt to distributed components
### Changed identity validation to use [JWT](https://github.com/golang-jwt/jwt) for distributed components
### Both front-end and back-end use i18n (international and localized)
### Painting presentation only made a small part of the changes to the main changes are to adapt to the cell phone terminal
<hr>

## Front-end Screen Preview
![index_en.png](imgs/index_en.png)
&emsp;
![Login_en.png](imgs/Login_en.png)

## Line Screen Preview
![lineSettingMenu_en.png](imgs/lineSettingMenu_en.png)
&emsp;
![lineLocationsList_en.png](imgs/lineLocationsList_en.png)

![setLocationRange_en.png](imgs/setLocationRange_en.png)
&emsp;
![setNotionRange_en.png](imgs/setNotionRange_en.png)

![lineMessage_en.png](imgs/lineMessage_en.png)
![setNotionRangeMessage_en.png](imgs/setNotionRangeMessage_en.png)
