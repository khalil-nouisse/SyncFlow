const { Pool } = require('pg');

const pool = new Pool({
  user: process.env.DB_USER || 'valury',
  password: process.env.DB_PASSWORD || 'your_password',
  host: process.env.DB_HOST || 'localhost',
  port: process.env.DB_PORT || 5432,
  database: process.env.DB_NAME || 'postgres'
});

module.exports = {
  query: (text, params) => pool.query(text, params)
};