const express = require('express');
const cors = require('cors');
const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());
app.use(cors({ origin: 'http://localhost:4200' })); // Angular origin

const grammarRoutes = require('./routes/grammarSch');
app.use('/grammarSch', grammarRoutes);

const plfsRoutes = require('./routes/plfs');
app.use('/plfs', plfsRoutes);

app.get('/', (req, res) => {
  res.send('MySQL API is running...');
});

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});
