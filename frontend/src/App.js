import logo from './logo.svg';
import './App.css';
import { Form, Input, Button, Checkbox, Card, Spin, message } from 'antd';
import {useState} from 'react'

import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

import ReactMarkdown from 'react-markdown'

import { Login} from './LogReg';

const { TextArea } = Input;

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

// just for test... remove this counter after test:
let counter = 50;


function ListOfNotes({updateNoteId, updateNoteState, token}) {
  const [notesList, updateNotesList] = useState([{id: 1}, {id: 123}]);

  const newNote = () => {
    // TODO: send POST request to the server
    counter += 1;
    const newID = counter;
    updateNotesList(x => [...x, {id: newID}]);
  };
  const selectNote = (id) => {
    // TODO: send GET request to server
    const success = true;
    if (success) {
      console.log("select", id);
      updateNoteState(x => NOTE_STATE.show); 
    }
  };
  const deleteNote = (id) => {
    // TODO: send DELETE request to server
    console.log("delete", id);
  };
  
  return (
    <>
      <Button type="primary" onClick={newNote}>ساخت یادداشت جدید</Button>
      <div dir="rtl">فهرست یادداشت‌ها:</div>
      <div>
      {notesList.map(function(d, idx){
        return (
          <>
            <div dir="rtl" key={idx}>
              {d.id}
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

function ShowNote({noteId, noteState, updateNoteState, token}) {
  const [editorText, updateEditorText] = useState("# Hello, *world*!\n## salam2\n### salam3\nhello ***majid***");
  
  const onChange = e => {
    updateEditorText(x => e.target.value);
  };

  const saveNote = () => {
    // TODO: send PUT request to server by noteID
  };

  const close = () => updateNoteState(x => NOTE_STATE.none);

  return (
    <>
      <TextArea showCount maxLength={1000} style={{ width: "90%", height: 200 }} onChange={onChange} defaultValue={editorText} />
      <div style={{ width: "90%", height: 200 }}>
        <ReactMarkdown>{editorText}</ReactMarkdown>
      </div>
      <div>
        <Button type="primary" onClick={saveNote}>ذخیره یادداشت</Button>
        <Button type="primary" onClick={close}>بستن یادداشت</Button>
      </div>
    </>
  )
}


function Dashboard({updateState, token, updateToken}) {
  const [noteId, updateNoteId] = useState(null);
  const [noteState, updateNoteState] = useState(NOTE_STATE.none);
  const log_out = () => {
    updateToken(x => null);
    updateState(x => STATES.login);
  };
  return (
    <>
      {
        <ListOfNotes
          updateNoteId={updateNoteId}
          updateNoteState={updateNoteState}
          token={token}
        />
      }
      {noteState === NOTE_STATE.show && <ShowNote token={token} noteId={noteId} noteState={noteState} updateNoteState={updateNoteState}/>}
      <div>
        <Button type="primary" onClick={log_out}>خروج از اکانت</Button>
      </div>
    </>
  )
}


function App() {
  const [state, updateState] = useState(STATES.login);
  const [token, updateToken] = useState(null);

  return (
    <>
      {state === STATES.login && <Login updateState={updateState} />}
      {state === STATES.dashboard && <Dashboard updateState={updateState} />}
      <ToastContainer />
    </>
  );
}

export default App;
