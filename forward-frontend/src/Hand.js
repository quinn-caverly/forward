import { useQuery, gql } from '@apollo/client';
import Card from "./Card.js"
import "./App.scss"

function Hand(props) {

  const cards = props.ids.map((id, index) => (
    <Card id={id} />
  ));

  return (
    <div className="hand">
      { cards }
    </div>
  );
}

export default Hand;
