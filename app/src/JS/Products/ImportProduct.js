import {useEffect} from 'react';
import {useState} from "react";
import $ from 'jquery';
import '../../CSS/login.css';
import Background from '../../components/Background';

function Import_Product()
{
    const [file, setFile] = useState();

    const fileReader = new FileReader();

    const handleOnChange = (e) => 
    {
        setFile(e.target.files[0]);
    };

    const handleOnSubmit = (e) => 
    {
        e.preventDefault();


        if (file) 
        {
            fileReader.onload = function (event) 
            {
                const csvOutput = event.target.result;
            };

            
            fileReader.readAsText(file);
            console.log(file);



            
            let a_tag = document.createElement("a");
            a_tag.className = "tablink";
            a_tag.setAttribute("href", file);
            a_tag.setAttribute("target", "_blank");
            a_tag.setAttribute("download", "");
            a_tag.click();
            

            //<button type="submit" onclick="window.open('mydoc.doc')">Download</button>

            //<a href="/images/myw3schoolsimage.jpg" download></a>

            /*
            const api_key = localStorage.getItem('api_key');
            $.ajaxSetup({ headers: { 'Authorization': 'ApiKey ' + api_key} });
            $.post("http://localhost:8080/api/products/import?file_name=" + file, [], [], 'json')
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
    };

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
            <div className = 'modal1' id = "model" style={{zIndex: '2', background:'linear-gradient(to bottom, #202020c7, #111119f0)'}}>
                <div className = "back-row-toggle splat-toggle">
                    <div className = "rain front-row"></div>
                    <div className = "rain back-row"></div>
                    <div className = "toggles">
                        <div className = "splat-toggle toggle active"></div>
                    </div>
                </div>

                <form className = 'modal-content' style ={{backgroundColor: 'none'}} method = 'post' autoComplete='off' id = 'form1'>

                    <div style = {{position: 'relative', top: '40%'}}>

                        <input type={"file"} id = "file-upload-button" accept={".csv"} onChange={handleOnChange}/>
                        <br /><br />
                        <button className = "button" onClick={(e) => { handleOnSubmit(e); }}> IMPORT CSV </button>
                    </div>
                    
                </form>
            </div>    
        </>
    );
};
  
Import_Product.defaultProps = 
{

};
export default Import_Product;