import logo from './logo.svg';
import './App.css';
import { Form, Input, Button, Checkbox, Card, Spin, message } from 'antd';
import {useState} from 'react'

import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

import ReactMarkdown from 'react-markdown'

import { Login, Register} from './LogReg';

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
  const [editorText, updateEditorText] = useState("# Hello, *world*!\n## salam2\n### salam3\nhello ***majid***");
  
  const onChange = e => {
    updateEditorText(x => e.target.value);
  };

  return (
    <>
      <TextArea showCount maxLength={1000} style={{ width: "90%", height: 250 }} onChange={onChange} defaultValue={editorText} />
      <div style={{ width: "90%", height: 250 }}>
        <ReactMarkdown>{editorText}</ReactMarkdown>
      </div>
    </>
  )
}


function Dashboard({updateState}) {
  /// ListOfNotes
  /// ShowNote if a note selected
  /// EditOrPostNote if edit or post note
  const [noteId, updateNoteId] = useState(null);
  const [noteState, updateNoteState] = useState(NOTE_STATE.show);
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
      <div>
        <Button type="primary" onClick={log_out}>خروج از اکانت</Button>
      </div>
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
