const tracer = require('./tracer')('node-service');

const express = require('express');
const axios = require('axios').default;

const app = express();
const PORT = 8081;


app.use(express.json());

app.get('/', async (req, res) => {
  // const response = await axios.get(`http://localhost:8083/`)
  response = ""
  res.status(201).send("hello from node\n" + response)
});

app.listen(PORT);
