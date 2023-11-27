import {useEffect} from 'react';
import Page2 from '../components/Page2';

import Detailed_Warehousing from '../components/semi-components/Product/detailed_warehousing';
import MyImage from '../media/3ada.png';

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
            <Page2 />
            
        </>
    );
}

export default Settings;