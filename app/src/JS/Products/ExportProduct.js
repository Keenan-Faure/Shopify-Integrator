import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';
import '../../CSS/login.css';
import Background from '../../components/Background';

function Export_Product()
{
    const[inputs, setInputs] = useState({});

    const handleChange = (event) =>
    {
        const name = event.target.name;
        const value = event.target.value;
        setInputs(values => ({...values, [name]: value}))
    }
    const ExportProduct = (event) =>
    {
        event.preventDefault();
        console.log(inputs);

        /*
        $.post("http://localhost:8080/api/login", JSON.stringify(inputs),[], 'json')
        .done(function( _data) 
        {
            console.log(_data);
        })
        .fail( function(xhr) 
        {
            alert(xhr.responseText);
        });
        */
    }

    useEffect(() =>
    {
        /* Fix any incorrect elements */
        let navigation = document.getElementById("navbar");
        let modal = document.getElementById("model");
        modal.style.display = "block";
        navigation.style.animation = "MoveRight 1.2s ease";
        navigation.style.position = "fixed";
        navigation.style.left = "0%";
        navigation.style.width = "100%";

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


    }, []);

    return (
        <>
            <Background />
            <div className = 'modal1' id = "model">
                <div className = "back-row-toggle splat-toggle">
                    <div className = "rain front-row"></div>
                    <div className = "rain back-row"></div>
                    <div className = "toggles">
                        <div className = "splat-toggle toggle active"></div>
                    </div>
                </div>

                <form className = 'modal-content' method = 'post' onSubmit={(event) => ExportProduct(event)} autoComplete='off' id = 'form1'>
                    <div className = 'modal-container' id = "main">
                        <label style = {{fontSize: '18px'}}><b>Export Product</b></label>
                        <br /><br />
                        <div className = "holder">
                            <label><b>Product_Code</b></label>
                            <br />
                            <span><input type = 'text' placeholder = "Product grouping code" name = "username" value = {inputs.username || ""}  onChange = {handleChange} required></input></span>
                            <br /><br />

                            <label><b>Product_Title</b></label>
                            <br />
                            <span><input type = 'password' placeholder = "Title of product" name = "password" value = {inputs.password || ""} onChange = {handleChange} required></input></span>
                            <br /><br />

                            <label><b>Product_Description</b></label>
                            <br />
                            <span><input type = 'text' placeholder = "Description of product" name = "username" value = {inputs.username || ""}  onChange = {handleChange} required></input></span>
                            <br /><br />
                        </div>

                        <button className = 'button' type = 'submit'>Add</button>
                    </div>
                </form>
            </div>    
        </>
    );
};
  
Export_Product.defaultProps = 
{

};
export default Export_Product;