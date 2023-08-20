import './App.scss';
import './styles.scss'
import Card from "./Card.js"
import Hand from "./Hand.js"
import Deck from "./Deck.js"

function App() {

  return (
    <div>
      <div class="main-header">
        <a href="https://github.com/quinn-caverly" target="_blank" class="main-title">
           github.com/quinn-caverly
        </a>
      </div>
      <div class="main">
        <Deck />
      </div>
    </div>
  );
}

export default App;
