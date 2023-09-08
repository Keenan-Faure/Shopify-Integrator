import {useEffect} from 'react';
import Page1 from '../components/Page1';
import '../CSS/page1.css';

/* Must start with a Caps letter */
function Orders(props)
{
    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        let model = document.getElementById("model");
        let main = document.querySelector(".main");
        window.onload = function(event)
        {
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
            model.style.animation = "none";
            main.style.animation = "SlideUp3 1.2s ease-in";
        }
        
    }, []);

    return (
        <>
            <Page1 
            filter_display = "none" main_bgc = "transparent" main_top = "53%" main_left = "50%" transform = "translate(-50%, -50%)"
            width = "70%" height = "90%" animation = "SlideUp3 1.2s ease-in"/>
            
        </>
    );
}

export default Orders;