import {useEffect} from 'react';
import '../CSS/dashboard.css';

function Dashboard(props)
{
    useEffect(()=> 
    {
        /* Ensures the navbar + model is set correctly */
        let navigation = document.getElementById("navbar");
        let model = document.getElementById("model");
        window.onload = function(event)
        {
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
            model.style.animation = "none";
            model.style.display = "none";
        }

        /* logout */
        let logout = document.getElementById("logout");
        logout.addEventListener("click", () =>
        {
            logout.style.display = "none";
            model.style.animation = "FadeIn ease-in 1s";
            model.style.display = "block";
        });
        
    }, []);

    return (
        <div className = "dashboard" id = "dashboard">
            <div className = "logout" id = "logout">Logout</div>
        </div>
    );
}


export default Dashboard;