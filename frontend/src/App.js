import logo from './logo.svg';
import './App.css';
import { Form, Input, Button, Checkbox, Card, Spin, message } from 'antd';
import {useState} from 'react'
/// use toast for error handeling

const STATES = {
  'register': 1,
  'login': 2,
  'dashboard': 3,
};

const NOTE_STATE = {
  'none': 0,
  'show': 1,
  'edit': 2,
  'post': 3,
};

function Register({updateState}) {
  console.log("Register");
  
  const register_request = () => {
    console.log("Sending register request to server...");
    // TODO: send request to server
    let success = true;
    if (success) {
      // toast
      console.log("ثبت نام با موفقیت انجام شد.");
      updateState(x => (STATES.login));
    } else {
      console.log("مشکلی در ثبت نام شما پیش آمد. دوباره تلاش کنید.");
    }
  };
  return (
    <>
      <div dir="rtl">سلام. لطفا ثبت نام کنید.</div>
      <Input placeholder="Username" />
      <Input placeholder="Password" />
      {/* TODO: Password holder must be like **** */}
      <Button type="primary" onClick={register_request}>ثبت نام</Button>
    </>
  )
}

function Login({updateState}) {
  /// username
  /// password
  /// link to go to register
  const login_request = () => {
    console.log("Sending login request to server...");
    // TODO: send login to server
    let success = true;
    if (success) {
      // toast
      console.log("ورودت موفقیت‌آمیز بود عزیزم!");
      updateState(x => (STATES.dashboard));
    } else {
      console.log("مشکلی در ورود شما پیش آمد. دوباره تلاش کنید.");
    }
  };
  const go_to_register = () => {
    // change state to registery
    updateState(x => (STATES.register));
  };
  return (
    <>
      <Input placeholder="Username" />
      <Input placeholder="Password" />
      <Button type="primary" onClick={login_request}>ورود</Button>
      <Button type="primary" onClick={go_to_register}>می‌خواهم اکانت بسازم</Button>
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

function Dashboard({updateState}) {
  /// ListOfNotes
  /// ShowNote if a note selected
  /// EditOrPostNote if edit or post note
  const [noteId, updateNoteId] = useState(null);
  const [noteState, updateNoteState] = useState(NOTE_STATE.none);
  const onSend = () => {};
  const log_out = () => {};
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
      <Button type="primary" onClick={log_out}>خروج از اکانت</Button>
    </>
  )
}


function App() {
  const [state, updateState] = useState(STATES.login);

  return (
    <>
      {state === STATES.login && <Login updateState={updateState} />}
      {state === STATES.register && <Register updateState={updateState} />}
      {state === STATES.dashboard && <Dashboard updateState={updateState} />}
    </>
  );
}

export default App;
