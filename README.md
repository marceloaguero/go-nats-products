# go-nats-products
Small example of a microservices application using NATS to communicate and a gateway to expose REST services

Este es un ejemplo de una mini aplicación que implementa un ABM de productos con un microservicio hecho en golang.

Provee una interfaz (API) REST hacia el "exterior", pero la comunicación "interna" utiliza una cola NATS (https://nats.io).
