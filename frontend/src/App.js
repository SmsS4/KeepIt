import logo from './logo.svg';
import './App.css';
import { Form, Input, Button, Checkbox, Card, Spin, message } from 'antd';
import {useState} from 'react'

import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

import ReactMarkdown from 'react-markdown'


const notify = (msg) => toast(msg);

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
    notify("Sending register request to server...");
    // TODO: send request to server
    let success = true;
    if (success) {
      notify("ثبت نام با موفقیت انجام شد.");
      updateState(x => (STATES.login));
    } else {
      notify("مشکلی در ثبت نام شما پیش آمد. دوباره تلاش کنید.");
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
    notify("Sending login request to server...");
    // TODO: send login to server
    let success = true;
    if (success) {
      notify("ورودت موفقیت‌آمیز بود عزیزم!");
      updateState(x => (STATES.dashboard));
    } else {
      notify("مشکلی در ورودت پیش اومد. دوباره تلاش کن.");
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
  // TODO:
  const selectNote = () => {
    
  };
  const deleteNote = () => {

  };
  return (
    <>
      <Button type="primary" onClick={selectNote}>انتخاب</Button>
      <Button type="primary" onClick={deleteNote}>پاک کردن</Button>
    </>
  )
}

function ShowNote({noteId, noteState, updateNoteState}) {
  /// markdown support
  /// save button
  const test_str = "# Hello, *world*!";
  return (
    <>
      <ReactMarkdown>{test_str}</ReactMarkdown>
    </>
  )
}

function EditOrPostNote({onSend, noteId, noteState, updateNoteState}) {
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
  const log_out = () => {
    notify("send log out request to server...");
    // TODO: send log out request to server
    let success = true;
    if (success) {
      notify("خروجت موفقیت‌آمیز بود عزیزم!");
      updateState(x => (STATES.login));
    } else {
      notify("مشکلی در خروجت پیش اومد. دوباره تلاش کن.");
    }
  };
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
          onSend={onSend}
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
      <ToastContainer />
    </>
  );
}

export default App;
