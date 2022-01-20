import logo from './logo.svg';
import './App.css';
import { Form, Input, Button, Checkbox, Card, Spin, message } from 'antd';
import {useState} from 'react'
/// use toast for error handeling

const STATES = {
  'register': 1,
  'login': 2,
  'dashboard': 3,
}

const NOTE_STATE = {
  'none': 0,
  'show': 1,
  'edit': 2,
  'post': 3,
}

function Register() {
  console.log("Register");
  
  const register_request = () => {
    console.log("Sending register request to server...")
    // TODO: send request to server
  };
  return (
    <>
      <div dir="rtl">سلام. لطفا ثبت نام کنید.</div>
      <Input placeholder="Username" />
      <Input placeholder="Password" />
      <Button type="primary" onClick={register_request}>ثبت نام</Button>
    </>
  )
}

function Login() {
  /// username
  /// password
  /// link to go to register
  return (
    <>

    </>
  )
}

function ListOfNotes({updateNoteId, updateNoteState}) {
  /// list of notes title
  /// delete button ?
  /// edit button ?
  /// show button
  return (
    <>
      
    </>
  )
}

function ShowNote({noteId, noteState}) {
  /// markdown support
  /// close button
  /// delete button ?
  /// edit button ?
  return (
    <>
      
    </>
  )
}

function EditOrPostNote({onSend}) {
  /// component for editing or posting note
  /// send button
  return (
    <>
      
    </>
  )
}

function Dashboard() {
  /// ListOfNotes
  /// ShowNote if a note selected
  /// EditOrPostNote if edit or post note
  const [noteId, updateNoteId] =  useState(null);
  const [noteState, updateNoteState] = useState(NOTE_STATE.none)
  const onSend = () => {};
  return (
    <>
      {
        <ListOfNotes
          updateNoteId={updateNoteId}
          updateNoteState={updateNoteState}
        />
      }
      {noteState === NOTE_STATE.show && <ShowNote noteId={noteId}  noteState={noteState} updateNoteState={updateNoteState}/>}
      {
        (noteState === NOTE_STATE.edit || noteState === NOTE_STATE.post) &&
        <EditOrPostNote 
          noteId={noteId}
          noteState={noteState}
          updateNoteState={updateNoteState}
        />
      }
    </>
  )
}


function App() {
  const [state, updateState] = useState(1);

  return (
    <>
      {state === STATES.login && <Login />}
      {state === STATES.register && <Register />}
      {state === STATES.dashboard && <Dashboard />}
    </>
  );
}

export default App;
