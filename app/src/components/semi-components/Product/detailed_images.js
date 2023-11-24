import {useEffect} from 'react';
import '../../../CSS/detailed.css';

function Detailed_Images(props)
{
    useEffect(()=> 
    {
        
        
        function showSlides() 
        {
            let slides = document.getElementsByClassName("mySlides");
            
            
            setTimeout(() =>
            {
                console.log(slides);
                
                if(slides.length <= 2)
                {
                    
                    for(let i = 0; i < slides.length; i++)
                    {
                        slides[i].style.display = "block"
                        slides[i].style.animationName = "";
                    }
                }
                else if(slides.length > 2)
                {
                    for (let i = 0; i < slides.length; i++) 
                    {
                        slides[i].style.display = "none";  
                    }
                    slideIndex++;
                    if (slideIndex > slides.length) {slideIndex = 1}    
                    slides[slideIndex-1].style.display = "block";  
                    setTimeout(() =>
                    {
                        showSlides();
                    }, 5000); // Change image every 5 seconds
                }
                else
                {
                    console.log("no slides to display");
                }
                
            },100)
            
            
            
        }

        let slideIndex = 0;
        showSlides();

    }, []);

    return (
    <>
        <div className = "mySlides fade">
            <img src = {props.Image1} className = "details-image"></img>
        </div>
    </>
    );
  
};

Detailed_Images.defaultProps = 
{ 
    Image1: '#ccc',
}
export default Detailed_Images;