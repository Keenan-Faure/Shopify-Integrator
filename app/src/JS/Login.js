import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';

import Background from "../components/Background";
import '../CSS/login.css';

function Login()
{
    const[inputs, setInputs] = useState({});
    const [result, setResult] = useState("");
    const [result2, setResult2] = useState("");
    const [result3, setResult3] = useState("");

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }

    const Login = (event) =>
    {
        event.preventDefault();
        const form = $(event.target);
        
        $.ajax
        ({
            type: "POST",
            url: "http://localhost:8000/api/login",
            data: form.serialize(), // important to maintain the form data output
            dataType: 'json',
            success(data) 
            {
                setResult3(data);
                console.log(data);
            },
        });
        
        
    }

    const Register = (event) =>
    {
        event.preventDefault();
        const form = $(event.target);
        /* Show the user the information received from the api */

        let re = document.querySelector(".result-container");
        re.style.display = "block";

        console.log(inputs.authentication);
        
        
        $.ajax
        ({
            type: "POST",
            url: "http://localhost:8080/api/register",
            data: form.serialize(), // important to maintain the form data output
            dataType: 'json',
            success(data) 
            {
                setResult2(data);
                console.log(data);
            },
        });
        
        
        /* Upon sucessfully registering, return to the login form to login */
        /* Information is retrieved via the post request */
    }

    const Register_auth = (event) =>
    {
        event.preventDefault();

        console.log(inputs);
        /*
        $.ajax
        ({
            type: "POST",
            url: "http://localhost:8080/api/preregister",
            data: 
            {
                name: inputs.name, 
                email: inputs.email
            },
            dataType: 'json',
            success(data) 
            {
                setResult(data);
                console.log(data);
            },
        });
        */
    
        $.post("http://localhost:8080/api/preregister", JSON.stringify(inputs))
        .done(function( _data) 
        {
            console.log(_data);
        })
        .fail( function(xhr, textStatus, errorThrown) 
        { 
            console.log(errorThrown.responseText);
            console.log(xhr.responseText)
            console.log(textStatus.responseText)
        });
        
        

        /* Adds the proceedings after the user enters his registration info */
        let form2 = document.getElementById("form2");
        let form3 = document.getElementById("form3");
        let return_button2 = document.querySelector(".return-button2");

        form2.style.animation = "Fadeout ease-out 1s";
        form2.style.display = "none";
        form3.style.animation = "FadeIn ease-in 1s";
        form3.style.display = "block";
        return_button2.style.display = "block";
        
        
        
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
                <div className = 'modal-container'>

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
                    
                    <label style = {{fontSize: '18px'}}><b><u>Register a New Account</u></b></label>
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
                <div className = 'modal-container'>
                    <div className = 'reg-portion' id = "reg-portion">
                        <label><b>Authentication Token</b></label>
                        <br /><br /><br />
                        <span><input type = 'password' placeholder = "Enter Token" name = "authentication" value = {inputs.authentication || ""} onChange = {handleChange} required></input></span>
                        <br /><br /><br />
                        <button className = 'button' id = "reg-auth" type = 'submit'>Register</button>
                    </div>
                </div>
            </form>

            <div className = 'result-container'>
                <div className = 'reg-portion'>
                    <label><b>Information Returned</b></label>
                    <br /><br /><br />
                    <span>Api Key</span>
                    <br /><br /><br />
                    <button className = 'button'type = 'button'>Copy to Clipboard</button>
                </div>
            </div>

            <div className = 'return-button'></div>
            <div className = 'return-button2'></div>
        </div>

        
    </div>
    </>
    );
  
};
    
export default Login;