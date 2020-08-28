import React from 'react';
import { BrowserRouter, Route } from 'react-router-dom'
import CreatePoll from './components/CreatePoll/CreatePoll';
import PollPage from './components/PollPage/PollPage';
import './App.css'

function App() {
  return (
    <BrowserRouter>
      <Route path="/createPoll" component={CreatePoll}></Route>
      <Route path="/poll/:uuid" component={PollPage}></Route>
    </BrowserRouter>
  );
}

export default App;
