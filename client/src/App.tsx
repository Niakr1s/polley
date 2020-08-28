import React from 'react';
import { BrowserRouter, Route } from 'react-router-dom'
import CreatePoll from './components/CreatePoll/CreatePoll';
import PollPage from './components/PollPage/PollPage';
import PollsPage from './components/PollsPage/PollsPage';
import './App.css'

function App() {
  return (
    <BrowserRouter>
      <Route path="/createPoll" component={CreatePoll}></Route>
      <Route path="/poll/:uuid" component={PollPage}></Route>
      <Route exact path="/" component={PollsPage}></Route>
    </BrowserRouter>
  );
}

export default App;
