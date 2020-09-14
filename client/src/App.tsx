import React from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom'
import CreatePoll from './components/CreatePoll/CreatePoll';
import PollPage from './components/PollPage/PollPage';
import PollsPage from './components/PollsPage/PollsPage';
import './App.css'
import { FourThousandFour } from './components/FourTHousandFour/FourThousandFour';

function App() {
  return (
    <BrowserRouter>
      <Switch>
        <Route path="/createPoll" component={CreatePoll}></Route>
        <Route path="/poll/:uuid" component={PollPage}></Route>
        <Route exact path="/" component={PollsPage}></Route>
        <Route component={FourThousandFour}></Route>
      </Switch>
    </BrowserRouter>
  );
}

export default App;
