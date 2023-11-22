import {useEffect} from 'react';
import Page2 from '../components/Page2';
import Detailed_product from '../components/semi-components/detailed_product';
import MyImage from '../media/Screenshot.png';

import '../CSS/page2.css';

/* Must start with a Caps letter */
function Settings()
{
    useEffect(()=> 
    {
        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            navigation.style.left = "0%";
            navigation.style.position = "relative";
            navigation.style.width = "100%";
        }
        
    }, []);

    return (
        <>
            <Detailed_product Product_Image = {MyImage}/>
            
        </>
    );
}

export default Settings;