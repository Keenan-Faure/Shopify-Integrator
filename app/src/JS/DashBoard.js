import {useEffect} from 'react';
import '../CSS/dashboard.css';

function Dashboard(props)
{
    useEffect(()=> 
    {
        /* Ensures the navbar + model is set correctly */
        let navigation = document.getElementById("navbar");
        let model = document.getElementById("model");
        let logout = document.getElementById("logout");
        window.onload = function(event)
        {
            
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
            
            navigation.style.display = "block";
            logout.style.display = "block";
            //model.style.animation = "FadeIn ease-in 1s";
            model.style.display = "none";
            
        }

        /* On form submit re-shape elements 
        let form1 = document.getElementById("form1");
        form1.onsubmit = function(event)
        {
            model.style.animation = "Fadeout 1s ease-out";
            setTimeout(() =>
            {
                model.style.display = "block";
                navigation.style.left = "0%";
                navigation.style.position = "relative";
                navigation.style.width = "100%";
                logout.style.display = "block";
            }, 1000);
           
        }
        */

        /* logout */
        logout.addEventListener("click", () =>
        {
            logout.style.display = "none";
            model.style.animation = "FadeIn ease-in 1s";
            model.style.display = "block";
            /* Session Destroy */
        });
        
    }, []);

    return (
        <div className = "dashboard" id = "dashboard">
            <div className = "logout" id = "logout">Logout</div>
        </div>
    );
}


export default Dashboard;