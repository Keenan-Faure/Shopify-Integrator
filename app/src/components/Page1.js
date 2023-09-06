import {useEffect} from 'react';
import Background from './Background';

import '../CSS/page1.css';


function Page1(props)
{
    useEffect(()=> 
    {

        /* Ensures the navbar is set correctly */
        let navigation = document.getElementById("navbar");
        window.onload = function(event)
        {
            
            navigation.style.left = "30%";
            navigation.style.position = "absolute";
            navigation.style.width = "70%";
            navigation.style.animation = "MoveLeft 1.2s ease";
        }

        /* animation for the pan elements */
        let pan = document.querySelectorAll(".pan");
        setTimeout(() =>
        {
            for(let i = 0; i < pan.length; i ++)
            {
                pan[i].style.display = "block";
                pan[i].style.display = "appear 1s ease-in";
            }
        }, 1400);


    }, []);



    return (
        <>
            <Background />
            <div className = "filter">

            </div>
            <div className = "main">
                <div className = "main-elements">
                    <div className = "pan"></div>
                    <div className = "pan"></div>
                    <div className = "pan"></div>
                    <div className = "pan"></div>
                </div>
            </div>
        </>
    );
}

export default Page1;