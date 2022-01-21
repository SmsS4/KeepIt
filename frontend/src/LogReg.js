import { ToastContainer, toast } from 'react-toastify';
import { Form, Input, Button, Checkbox, Card, Spin, message } from 'antd';


const notify = (msg) => toast(msg);

const STATES = {
  'register': 1,
  'login': 2,
  'dashboard': 3,
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


export {Register, Login}