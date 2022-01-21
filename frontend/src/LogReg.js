import { ToastContainer, toast } from 'react-toastify';
import {  Input, Button, Checkbox, Card, Spin, message } from 'antd';
import React from 'react';
import ReactDOM from 'react-dom';

const BASE = "http://localhost:8080"
const LOGIN_URL = `${BASE}/pub/login/`
const REGISTER_URL = `${BASE}/pub/register/`

const notify = (msg) => toast(msg);

const STATES = {
  'login': 2,
  'dashboard': 3,
};


const appStyle = {
  display: 'flex',
  height: '980px',
  backgroundImage: "url(back.jpeg)",
  backgroundSize: '100% 100%',
  backgroundRepeat: 'no-repeat',
};

const formStyle = {
  transform: "scale(2, 2)",
  marginLeft: '40%',
  marginTop: '28%',
  padding: '10px',
  // border: '1px solid #c9c9c9',
  // borderRadius: '5px',
  // background: '#f5f5f5',
  width: '220px',
  display: 'block'
};



const labelStyle = {
  margin: '10px 0 5px 0',
  fontFamily: 'Arial, Helvetica, sans-serif',
  fontSize: '15px',
};

const inputStyle = {
  margin: '5px 0 10px 0',
  padding: '5px', 
  border: '1px solid #bfbfbf',
  borderRadius: '3px',
  boxSizing: 'border-box',
  width: '100%'
};

const submitStyle = {
  margin: '10px 0 0 0',
  padding: '7px 10px',
  border: '1px solid #efffff',
  borderRadius: '3px',
  background: '#3085d6',
  width: '100%', 
  fontSize: '15px',
  color: 'white',
  display: 'block'
};

const Field = React.forwardRef(({label, type}, ref) => {
  return (
    <div>
      <label style={labelStyle} >{label}</label>
      <input ref={ref} type={type} style={inputStyle} />
    </div>
  );
});

const Form = ({onLogin, onRegister}) => {
  const usernameRef = React.useRef();
  const passwordRef = React.useRef();
  var [isLogin, setIsLogin] = React.useState(false);
  const handleSubmit = e => {
      e.preventDefault();
      const data = {
          username: usernameRef.current.value,
          password: passwordRef.current.value
      };
      if (isLogin) {
        onLogin(data);
      } else {
        onRegister(data)
      }
  };
  return (
    <form style={formStyle} onSubmit={handleSubmit} >
      <Field ref={usernameRef} label="Username:" type="text" />
      <Field ref={passwordRef} label="Password:" type="password" />
      <div>
        <button onClick={() => (setIsLogin(true))} style={submitStyle} type="submit">{"ورود"}</button>
        <button onClick={() => (setIsLogin(false))} style={submitStyle} type="submit">{"ثبت نام"}</button>
      </div>
    </form>
  );
};


function Login({updateState}) {
    /// username
    /// password
    /// link to go to register
    const login_request = (data) => {
      notify("Sending login request to server...");
      let success = true;
      if (success) {
        notify("ورودت موفقیت‌آمیز بود عزیزم!");
        updateState(x => (STATES.dashboard));
      } else {
        notify("مشکلی در ورودت پیش اومد. دوباره تلاش کن.");
      }
    };

    const register_request = (data) => {
      notify("Sending register request to server...");

      fetch(REGISTER_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      })
      .then(response => response.json())
      .then(data => {
        console.log('Success:', data);
      })
      .catch((error) => {
        console.error('Error:', error);
      });

      let success = true;
      if (success) {
        notify("ثبت نام با موفقیت انجام شد.");
      } else {
        notify("مشکلی در ثبت نام شما پیش آمد. دوباره تلاش کنید.");
      }
    };

    return (
      <>
      <div style={appStyle}>
        <Form onRegister={register_request} onLogin={login_request} reg={false}/>
      </div>
      </>
    )
  }


export {Login}