# 智慧型空氣品質監測系統 v2
[English](README.md) | 繁體中文

<a href="https://github.com/wuchieh"><img src="imgs/LogoWhile.png" alt="wuchieh" style="height: 100px;"></a>
&emsp;
<a href="https://github.com/MeowXiaoXiang"><img src="https://github.com/MeowXiaoXiang.png" alt="MeowXiaoXiang" style="height: 100px;"></a>

## 需求
```
PostgreSQL
Redis
go version > 1.20
```

## 介紹
- 此為 `智慧型空氣品質監測系統` 的第二代
- 優化了絕大多數的系統
- 在新版本中 引進了前後端分離 
- 前端使用 Vue3 框架 而非全交由後端渲染 以降低後端壓力
- 並且更換了數據庫 由 MySQL 改成了 [PostgreSQL](https://github.com/lib/pq)
- 以及引入了 [Redis](https://github.com/redis/go-redis/) 以提升部分效率 及適配分布式部屬
- 為了分布式部屬 身分驗證更改為使用 [JWT](https://github.com/golang-jwt/jwt)
- 前後端同時都有使用 i18n(國際化與在地化)
- 畫面呈現上只做了小部分的改動 主要的修改都是為了適配手機端

<hr>

## 前端畫面預覽
![index_tw.png](imgs/index_tw.png)
&emsp;
![Login_tw.png](imgs/Login_tw.png)

## Line畫面預覽
![lineSettingMenu_tw.png](imgs/lineSettingMenu_tw.png)
&emsp;
![lineLocationsList_tw.png](imgs/lineLocationsList_tw.png)

![setLocationRange_tw.png](imgs/setLocationRange_tw.png)
&emsp;
![setNotionRange_tw.png](imgs/setNotionRange_tw.png)

![lineMessage_tw.png](imgs/lineMessage_tw.png)
![setNotionRangeMessage_tw.png](imgs/setNotionRangeMessage_tw.png)