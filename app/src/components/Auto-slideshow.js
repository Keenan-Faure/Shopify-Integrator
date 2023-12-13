import {useEffect} from 'react';
import '../CSS/Auto-slideshow.css';
import Image1 from '../media/Slide/1.png';
import Image2 from '../media/Slide/2.png';
import Image3 from '../media/Slide/3.png';


function Auto_Slideshow(props)
{
    useEffect(()=> 
    {
        let slideIndex = 0;
        showSlides();
        
        function showSlides() 
        {
          let i;
          let slides = document.getElementsByClassName("myySlides");
          for (i = 0; i < slides.length; i++) 
          {
            slides[i].style.display = "none";  
          }
          slideIndex++;
          if (slideIndex > slides.length) {slideIndex = 1}    
          slides[slideIndex-1].style.display = "block";  
          setTimeout(() =>
            {
                showSlides();
            }, 5000); // Change image every 2 seconds
        }
    }, []);

    return (
    <>
        <div className = "auto-slideshow-container" style = {{display: props.Display}}>

            <div className = "myySlides fade">
                <img src = {Image1} style = {{width: '100%'}}></img>
            </div>

            <div className = "myySlides fade">
                <img src = {Image2} style = {{width: '100%'}}></img>
            </div>

            <div className = "myySlides fade">
                <img src = {Image3} style = {{width: '100%'}}></img>
            </div>
        </div>
    </>
    );
  
};

export default Auto_Slideshow;