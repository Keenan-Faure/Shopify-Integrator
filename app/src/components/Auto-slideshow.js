import {useEffect} from 'react';
import '../CSS/Auto-slideshow.css';


function Auto_Slideshow(props)
{
    useEffect(()=> 
    {
        let slideIndex = 0;
        showSlides();
        
        function showSlides() 
        {
          let i;
          let slides = document.getElementsByClassName("mySlides");
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

            <div className = "mySlides fade">
                <img src = {props.Image1} style = {{width: '100%'}}></img>
            </div>

            <div className = "mySlides fade">
                <img src = {props.Image2} style = {{width: '100%'}}></img>
            </div>

            <div className = "mySlides fade">
                <img src = {props.Image3} style = {{width: '100%'}}></img>
            </div>
        </div>
    </>
    );
  
};

Auto_Slideshow.defaultProps = 
{ 
    Image1: '#ccc',
    Image2: '#ccc',
    Image3: '#ccc',
}
export default Auto_Slideshow;