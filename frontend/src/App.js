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
  // new note button
  const [notesList, updateNotesList] = useState([]);
  const get_notes_list = () => {
    // TODO: sending request to get list of notes
    const notesListRes = [
      {id: 1, title:"first title"},
      {id: 123, title:"second"}
    ];
    updateNotesList(x => notesListRes);
  };

  const newNote = () => {
    // TODO:
  };
  const selectNote = (id) => {
    // TODO:
    console.log("select", id);
  };
  const deleteNote = (id) => {
    // TODO: send delete request to server
    console.log("delete", id);
  };
  
  return (
    <>
      <Button type="primary" onClick={get_notes_list}>به‌روزرسانی فهرست</Button>
      <div dir="rtl">فهرست یادداشت‌ها:</div>
      <div>
      {notesList.map(function(d, idx){
        return (
          <>
            <div dir="rtl" key={idx}>
              {d.title}
              <Button type="primary" onClick={() => selectNote(d.id)}>انتخاب</Button>
              <Button type="primary" onClick={() => deleteNote(d.id)}>پاک کردن</Button>
            </div>
          </>
        )
      })}
      </div>
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


function Dashboard({updateState}) {
  /// ListOfNotes
  /// ShowNote if a note selected
  /// EditOrPostNote if edit or post note
  const [noteId, updateNoteId] = useState(null);
  const [noteState, updateNoteState] = useState(NOTE_STATE.none);
  const onSend = () => {};
  const log_out = () => updateState(x => (STATES.login));
  return (
    <>
      {
        <ListOfNotes
          updateNoteId={updateNoteId}
          updateNoteState={updateNoteState}
        />
      }
      {noteState === NOTE_STATE.show && <ShowNote noteId={noteId}  noteState={noteState} updateNoteState={updateNoteState}/>}
      <Button type="primary" onClick={log_out}>خروج از اکانت</Button>
    </>
  )
}


function App() {
  const [state, updateState] = useState(STATES.login);
  // TODO: token, updateToken

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
