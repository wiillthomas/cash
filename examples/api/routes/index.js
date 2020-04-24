var express = require('express');
var router = express.Router();

const axios = require("axios")

/* Add to cache */
router.get('/addBulk/:value', function(req, res, next) {
  const { value } = req.params;
  
  const times = []

  const num = 1000

  let gotRequests = 0

  for ( let i = 0; i < num; i += 1 ) {

    let end
    setTimeout(() => {
      let start = new Date()
      axios.get(`http://cache:9192/create?value=${ Math.random() }`)
        .then((response) => {
          end = new Date() 
          total = end.getTime() - start.getTime()
          times.push(`${response.data} ${total}`)
          console.log(total)
          gotRequests += 1
          if ( gotRequests == num ) {
            res.status(200).send(times)
          }
        })
        .catch((err) => {
          // console.log(err)
        })
    }, i * 2 )
  }
});
 
/* Bulk add to cache */
router.get('/add/:value', function(req, res, next) {
  const { value } = req.params;

  let start = new Date()
  let end
  axios.get(`http://cache:9192/create?value=${ value }&?expiry=7000`)
    .then((response) => {
      end = new Date()
      total = end.getTime() - start.getTime()
      console.log(total)
      res.status(200).send(`hash: ${response.data} op: ${ total }`)
    })
  .catch((err) => console.log(err))
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

module.exports = router;
