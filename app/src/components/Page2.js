import {useEffect} from 'react';
import Background from './Background';

import '../CSS/page2.css';


function Page2(props)
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
    }, []);

    return (
        <>
            <Background />
            <div className = "component1">
                <div className = "main-container">
                    <div className = "settings">

                        <div className = "settings-section 1">
                            <div className = "settings-title">Setting 1</div>
                            <div className = "settings-description">Description goes here of this setting</div>

                            <button className = "settings-button on">On</button>
                            <button className = "settings-button off">Off</button>
                        </div>

                        <div className = "settings-section 2">
                            <div className = "settings-title">Setting 1</div>
                            <div className = "settings-description">Description goes here of this setting</div>

                            <button className = "settings-button on">On</button>
                            <button className = "settings-button off">Off</button>
                        </div>

                        <div className = "settings-section 3">
                            <div className = "settings-title">Setting 1</div>
                            <div className = "settings-description">Description goes here of this setting</div>

                            <button className = "settings-button on">On</button>
                            <button className = "settings-button off">Off</button>
                        </div>

                        <div className = "settings-section 3">
                            <div className = "settings-title">Setting 1</div>
                            <div className = "settings-description">Description goes here of this setting</div>

                            <button className = "settings-button on">On</button>
                            <button className = "settings-button off">Off</button>
                        </div>
                        
                         <div className = "center">
                            <div className = "pagination">
                                <a href = "#">&laquo;</a>
                                <a href = "#" className = "activee">1</a>
                                <a href = "#">2</a>
                                <a href = "#">3</a>
                                <a href = "#">4</a>
                                <a href = "#">5</a>
                                <a href = "#">6</a>
                                <a href = "#">&raquo;</a>
                            </div>
                         </div>
                    </div>
                </div>
                
            </div>
            
        </>
    );    

}    

export default Page2;