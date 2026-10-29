-- ========================================================
-- Script de inicialización para PostgreSQL en Docker
-- Se ejecuta automáticamente la primera vez que se inicia el contenedor
-- ========================================================

-- Extensiones útiles para UUIDs y funciones criptográficas
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Schema global requerido por la arquitectura multi-tenant de ClientHub
CREATE SCHEMA IF NOT EXISTS global;

