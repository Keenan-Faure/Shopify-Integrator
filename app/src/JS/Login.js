import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';

import Background from "../components/Background";
import '../CSS/login.css';

function Login()
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }

    const Login = (event) =>
    {
        event.preventDefault();

        /* Hide the user the information received from the api */
        let re = document.querySelector(".result-container");
        re.style.display = "none";

        console.log(inputs);

        let message = document.getElementById("message");
        message.style.display = "block";

        $.ajaxSetup
        ({
            headers: { 'Authorization': 'ApiKey ' + inputs.password }
        });

        $.post("http://localhost:8080/api/login", JSON.stringify(inputs),[], 'json')
        .done(function( _data) 
        {
            console.log(_data);

            /* Sets the user information for this session */
            localStorage.setItem('api_key', inputs.password);

            message.innerHTML = "Login Sucessful";
            message.style.backgroundColor = "rgb(6, 133, 6)";
            setTimeout(() =>
            {
                message.innerHTML = "";
                message.style.backgroundColor = "transparent";
                message.style.display = "none";
                window.location.href = '/dashboard';
            }, 1000);
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
            message.innerHTML = "Error - Login failed";
            message.style.backgroundColor = "rgb(175, 11, 11)";
            setTimeout(() =>
            {
                message.innerHTML = "";
                message.style.backgroundColor = "transparent";
                message.style.display = "none";
            }, 2000);
        });
    }

    const Register = (event) =>
    {
        event.preventDefault();

        /* Show the user the information received from the api */
        let re = document.querySelector(".result-container");
        let copyText = document.getElementById("myInput");
        let main2 = document.getElementById("main2");
        let message = document.getElementById("message");


        message.style.display = "block";
        
        $.post("http://localhost:8080/api/register", JSON.stringify(inputs),[], 'json')
        .done(function( _data) 
        {
            console.log(_data);
            copyText.innerHTML = JSON.stringify(_data, null, 2);

            message.innerHTML = "Registration Sucessful";
            message.style.backgroundColor = "rgb(6, 133, 6)";
            main2.style.animation = "Fadeout 2s ease-out";
            setTimeout(() =>
            {
                main2.style.display = "none";
                re.style.display = "block";
                message.innerHTML = "";
                message.style.backgroundColor = "transparent";
                message.style.display = "none";
            }, 2000);
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
            message.innerHTML = "Error - Ensure the Token is correct";
            message.style.backgroundColor = "rgb(175, 11, 11)";
            setTimeout(() =>
            {
                message.innerHTML = "";
                message.style.backgroundColor = "transparent";
                message.style.display = "none";
            }, 2000);
        });
    }

    const Register_auth = (event) =>
    {
        event.preventDefault();
        console.log(inputs);
        let message = document.getElementById("message");
        message.style.display = "block";

        $.post("http://localhost:8080/api/preregister", JSON.stringify(inputs),[], 'json')
        .done(function( _data) 
        {
            console.log(_data);

            message.innerHTML = "Email sent";
            message.style.backgroundColor = "rgb(6, 133, 6)";
            setTimeout(() =>
            {
                message.innerHTML = "";
                message.style.backgroundColor = "transparent";
                message.style.display = "none";
            }, 2000);

            /* Adds the proceedings after the user enters his registration info */
            let form2 = document.getElementById("form2");
            let form3 = document.getElementById("form3");
            let return_button2 = document.querySelector(".return-button2");
            form2.style.animation = "Fadeout ease-out 1s";
            form2.style.display = "none";
            form3.style.animation = "FadeIn ease-in 1s";
            form3.style.display = "block";
            return_button2.style.display = "block";

            
        })
        .fail( function(xhr) 
        { 
            alert(xhr.responseText);
            message.innerHTML = "Error";
            message.style.backgroundColor = "rgb(175, 11, 11)";
            setTimeout(() =>
            {
                message.innerHTML = "";
                message.style.backgroundColor = "transparent";
                message.style.display = "none";
            }, 2000);
        }); 
    }


    useEffect(()=> 
    {
        /* Ensure the model is shown */
        let model = document.getElementById("model");
        let navbar = document.getElementById("navbar");
        navbar.style.display = "none";
        model.style.display = "block";
        
        /* The swapping of forms */
        let register_button = document.getElementById("reg");
        let return_button = document.querySelector(".return-button");
        let return_button2 = document.querySelector(".return-button2");
        let form1 = document.getElementById("form1");
        let form2 = document.getElementById("form2");
        let form3 = document.getElementById("form3");
        register_button.addEventListener("click", () =>
        {
            form1.style.animation = "Fadeout ease-out 1s";
            form1.style.display = "none";

            form2.style.animation = "FadeIn ease-in 1s";
            form2.style.display = "block";

            return_button.style.display = "block";
        });

        /* return button swapping of forms */
        return_button.addEventListener("click", () =>
        {
            form2.style.animation = "Fadeout ease-out 1s";
            form2.style.display = "none";

            form1.style.animation = "FadeIn ease-in 1s";
            form1.style.display = "block";

            return_button.style.display = "none";

        });

        /* return button2 swapping of forms */
        return_button2.addEventListener("click", () =>
        {
            form3.style.animation = "Fadeout ease-out 1s";
            form3.style.display = "none";

            form2.style.animation = "FadeIn ease-in 1s";
            form2.style.display = "block";

            return_button2.style.display = "none";
            return_button.style.display = "block";
        });

        
        let clip = document.getElementById("clip");
        clip.addEventListener("click", () =>
        {
            let copyText = document.getElementById("myInput");
          
            navigator.clipboard.writeText(copyText.innerHTML);
            
            setTimeout(() =>
            {
                let re = document.querySelector(".result-container");
                let form2 = document.getElementById("form2");
                let message = document.getElementById("message");

                message.style.display = "block";
                message.innerHTML = "Copied to Clipboard";

                message.style.backgroundColor = "rgb(6, 133, 6)";
                re.style.animation = "Fadeout 2s ease-out";
                form2.style.animation = "Fadeout 2s ease-out";
                
                setTimeout(() =>
                {
                    window.location.reload();
                }, 1000);
            }, 1000)
            
        });



        /* Rain Functions */
        var makeItRain = function() 
        {
            //clear out everything
            $('.rain').empty();
          
            var increment = 0;
            var drops = "";
            var backDrops = "";
          
            while (increment < 100) 
            {

                //couple random numbers to use for various randomizations
                //random number between 98 and 1
                var randoHundo = (Math.floor(Math.random() * (98 - 1 + 1) + 1));
                //random number between 5 and 2
                var randoFiver = (Math.floor(Math.random() * (5 - 2 + 1) + 2));
                //increment
                increment += randoFiver;
                //add in a new raindrop with various randomizations to certain CSS properties
                drops += '<div class="drop" style="left: ' + increment + '%; bottom: ' 
                + (randoFiver + randoFiver - 1 + 100) + '%; animation-delay: 0.' + randoHundo 
                + 's; animation-duration: 0.5' + randoHundo + 's;"><div class="stem" style="animation-delay: 0.' 
                + randoHundo + 's; animation-duration: 0.5' + randoHundo + 's;"></div><div class="splat" style="animation-delay: 0.' 
                + randoHundo + 's; animation-duration: 0.5' + randoHundo + 's;"></div></div>';
                
                backDrops += '<div class="drop" style="right: ' + increment + '%; bottom: ' 
                + (randoFiver + randoFiver - 1 + 100) + '%; animation-delay: 0.' + randoHundo 
                + 's; animation-duration: 0.5' + randoHundo + 's;"><div class="stem" style="animation-delay: 0.' 
                + randoHundo + 's; animation-duration: 0.5' + randoHundo + 's;"></div><div class="splat" style="animation-delay: 0.' 
                + randoHundo + 's; animation-duration: 0.5' + randoHundo + 's;"></div></div>';
            }
          
            $('.rain.front-row').append(drops);
            $('.rain.back-row').append(backDrops);
        }
          
        $('.splat-toggle.toggle').on('click', function() 
        {
            $('body').toggleClass('splat-toggle');
            $('.splat-toggle.toggle').toggleClass('active');
            makeItRain();
        });
          
          $('.back-row-toggle.toggle').on('click', function() 
            {
                $('body').toggleClass('back-row-toggle');
                $('.back-row-toggle.toggle').toggleClass('active');
                makeItRain();
            });
          
            $('.single-toggle.toggle').on('click', function() 
            {
                $('body').toggleClass('single-toggle');
                $('.single-toggle.toggle').toggleClass('active');
                makeItRain();
            });
          
        /* DONT MAKE IT RAIN YET! */
        //makeItRain();

        /* Upon a successful login 
        function sucessfulLogin()
        {
            let model = document.getElementById("model");
            model.style.zIndex = "0";
        }
        */

    }, []);


    return (
    <>
    <Background />
    <div>
        <div className = 'modal1' id = "model" style = {{display: 'block'}}>
            <div className = "back-row-toggle splat-toggle">
                <div className = "rain front-row"></div>
                <div className = "rain back-row"></div>
                <div className = "toggles">
                    <div className = "splat-toggle toggle active"></div>
                </div>
            </div>

            <form className = 'modal-content' method = 'post' onSubmit={(event) => Login(event)} autoComplete='off' id = 'form1'>
                <div className = 'modal-container' id = "main">

                    <label style = {{fontSize: '18px'}}><b>Welcome. Please login to proceed</b></label>
                    <br /><br /><br />
                    <label><b>Username</b></label>
                    <br />
                    <span><input type = 'text' placeholder = "Name" name = "username" value = {inputs.username || ""}  onChange = {handleChange} required></input></span>
                    <br /><br /><br />
                    <label><b>Api Key</b></label>
                    <br />
                    <span><input type = 'password' placeholder = "Api-Key" name = "password" value = {inputs.password || ""} onChange = {handleChange} required></input></span>
                    <br /><br />
                    <button className = 'button' type = 'submit'>Proceed</button> <div id = "reg" className = 'text'>Or Register</div>
                </div>
            </form>

            <form style = {{display: 'none'}} className = 'modal-content'  method = 'post' onSubmit={(event) => Register_auth(event)} autoComplete='off' id = 'form2'>
                <div className = 'modal-container'>
                    
                    <label id = "info" style = {{fontSize: '18px'}}><b><u>Register a New Account</u></b></label>
                    <br /><br /><br />
                    <label><b>Username</b></label>
                    <br />
                    <span><input type = 'text' placeholder = "Name" name = "name" value = {inputs.name || ""}  onChange = {handleChange} required></input></span>
                    <br /><br /><br />
                    <label><b>Email</b></label>
                    <br />
                    <span><input type = 'email' placeholder = "Email" name = "email" value = {inputs.email || ""} onChange = {handleChange} required></input></span>
                    <br /><br />
                    <button id = "proc" className = 'button' type = 'submit'>Send Token</button>
                    <br /><br />
                </div>
            </form>

            <form style = {{display: 'none'}} className = 'modal-content'  method = 'post' onSubmit={(event) => Register(event)} autoComplete='off' id = 'form3'>
                <div className = 'modal-container' id = "main2">
                    <div className = 'reg-portion' id = "reg-portion">
                        <label><b>Authentication Token</b></label>
                        <br />
                        <div className = "message">A Token was sent to the email address</div>
                        <br /><br />
                        <span><input type = 'password' placeholder = "Enter Token" name = "token" value = {inputs.token || ""} onChange = {handleChange} required></input></span>
                        <br /><br /><br />
                        <button className = 'button' id = "reg-auth" type = 'submit'>Register</button>
                    </div>
                </div>
            </form>
            <div className = 'result-container'>
                <div className = 'reg-portion'>
                    <label><b>Information Returned</b></label>
                    <div className = "message">You are recommended to save this information!</div>
                    <br />
                    <pre id = "myInput" className = "pre"/>
                    <br /><br />
                    <button className = 'button' type = 'button' id = "clip">Copy to Clipboard</button>
                </div>
            </div>

            

            <div className = 'return-button'></div>
            <div className = 'return-button2'></div>
            <div className = 'info-message' id = 'message'></div>
        </div>

        
    </div>
    </>
    );
  
};
    
export default Login;