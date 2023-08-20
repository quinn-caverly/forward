import { useQuery, gql } from '@apollo/client';
import Hand from "./Hand.js"
import "./App.scss"
import React, { useState } from 'react';

import { IconButton } from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';

const COUNT_USABLE_PRODUCTS = gql`
  {
    countUsableProducts
  }
`;

const GENERATE_RANDOM_PRODS = gql`
  query GenerateRandomProds($Num: Int!) {
    generateRandomProds(Num: $Num)
  }
`;

function Deck() {
  const { loading, error, data } = useQuery(GENERATE_RANDOM_PRODS, {
    variables: { Num: 100 },
  });

  const [start, setStartIndex] = useState(0);

  const handleMoreButton = (valToAdd) => () => {
      setStartIndex(start + valToAdd)
  }

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error: {error.message}</p>;

  return (
    <div className="deck">
      <IconButton className="left-button" onClick={handleMoreButton(-10)}>
        <ArrowBackIcon fontSize="large"/>
      </IconButton>
      <Hand ids={data.generateRandomProds.slice(start, start+10)} />
      <IconButton className="right-button">
        <ArrowForwardIcon fontSize="large" onClick={handleMoreButton(10)}/>
      </IconButton>
    </div>
  );
}

// <Hand ids={data.GenerateRandomProds.generateRandomProds} />

export default Deck;
