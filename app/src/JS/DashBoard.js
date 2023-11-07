import {useEffect} from 'react';
import '../CSS/dashboard.css';

function Dashboard(props)
{
    useEffect(()=> 
    {
        /* Ensures the navbar + model is set correctly */
        let navigation = document.getElementById("navbar");
        let logout = document.getElementById("logout");
        window.onload = function(event)
        {
            
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
            navigation.style.display = "block";
            logout.style.display = "block"; 
        }

        /* logout */
        logout.addEventListener("click", () =>
        {
            logout.style.display = "none";
            navigation.style.display = "none";

            /* Session Destroy */
            window.location.href = '/';
        });
        
    }, []);

    return (
        <div className = "dashboard" id = "dashboard">
            <div className = "container">
                <div className = "logout" id = "logout">Logout</div>
            </div>
        </div>
    );
}


export default Dashboard;