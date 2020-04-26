var express = require('express');
var router = express.Router();

require("http").globalAgent.maxFreeSockets = 1000000
require('http').globalAgent.maxSockets = 1000000
const http = require("http")

console.log(require('http').globalAgent.maxSockets)

const axios = require("axios")
 
/* Benchmark Bulk add to cache in 1s */
router.get('/addBulk/:requests', function(req, res, next) {
  const { requests=100 } = req.params;
  
  const times = []
  
  const interval = requests / 1000

  let gotRequests = 0

  for ( let i = 0; i < requests; i += 1 ) {

    let end
    setTimeout(() => {
      let start = new Date()
      axios.get(`http://cache:9192/create?value=testval`, { httpAgent: new http.Agent({ keepAlive: true })})
        .then((response) => {
          end = new Date() 
          total = end.getTime() - start.getTime()
          times.push(total)
          console.log(total)
          gotRequests += 1
          if ( gotRequests == requests ) {
            let avgTime = 0;
            for ( let i = 0; i < times.length; i ++ ) {
                avgTime += times[i]
            }
            avgTime = avgTime / times.length
            res.status(200).send(`time: ${times} 
            avg: ${ avgTime}`)
          }
        })
        .catch((err) => {
          // console.log(err)
        })
    }, interval  )
  }
});
 
/* Add Single To Cache */
router.get('/add/:value/:expiry', function(req, res, next) {
  const { value, expiry } = req.params;

  let start = new Date()
  let end
  axios.get(`http://cache:9192/create?value=${ value }&?expiry=${ expiry }`)
    .then((response) => {
      end = new Date()
      total = end.getTime() - start.getTime()
      console.log(total)
      res.status(200).send(`hash: ${response.data} op: ${ total }`)
    })
  .catch((err) => console.log(err))
});


/* Continuously add items to cache */
router.get('/flood/', function(req, res, next) {
  setInterval(() => { 
    for ( let i = 0; i < 100 ; i += 1 ) {
        setTimeout(() => {
          axios.get(`http://cache:9192/create?value=testval`)
            .then((response) => {
            }) 
            .catch((err) => {
		console.log(err)
            })
        }, 10 )
    }
  }, 1000)
});


/* Read from cache */
router.get('/read/:hash', function(req, res, next) {
  const { hash } = req.params;

  let start = new Date()
  let end
  axios.get(`http://cache:9192/read?key=${ hash }`)
    .then((response) => {
      end = new Date()
      total = end.getTime() - start.getTime()
      console.log(total)
      res.status(200).send(`value: ${response.data} op: ${ total }`)
    })
  .catch((err) => console.log(err))
});

/* Delete from cache */
router.get('/destroy/:hash', function(req, res, next) {
  const { hash } = req.params;

  let start = new Date()
  let end
  axios.get(`http://cache:9192/destroy?key=${ hash }`)
    .then((response) => {
      end = new Date()
      total = end.getTime() - start.getTime()
      console.log(total)
      res.status(200).send(`value: ${response.data} op: ${ total }`)
    })
  .catch((err) => console.log(err))
});


/* Purge cache */
router.get('/purge', function(req, res, next) {
  const { hash } = req.params;

  let start = new Date()
  let end
  axios.get(`http://cache:9192/purge`)
    .then((response) => {
      end = new Date()
      total = end.getTime() - start.getTime()
      console.log(total)
      res.status(200).send(`value: ${response.data} op: ${ total }`)
    })
  .catch((err) => console.log(err))
});


module.exports = router;
