# PixelsZoom


    Pixel Zoom improve image resolution of an image; Image Resolution is determined by many factors such as Level of Zoom, Shading, Object detection and filters . Currently this package only support positive level of zoom; In Future , negative level of zoom will be implemented.



# LISTEN @ 

    Port = 3000
    Host = 'localhost'


# Routes 

    route_name                  method                       
      /                        get, post

# Arch 

 [Upload Image] -> [Zoom_in translator] == (Image)* -> [Zoom_out translator] == (Image)^

  * = Zoom_in , ^ = Zoom_out

  User will upload any image that will Zoom_in. Then translator process and return zoom in pixel image. By default reverse process is enabled so that Zoom_in pixel image will be reversed

# Data

    In images directory, myAvatar.png zoom in pixel image
    In images directory, myAvatar(1).png zoom out pixel image

    
# Test Data

    Open images/ myAvatar.png  in vscode or editor click touch pad or double click it will zoom without any blurnessness during scaling. 

    Open images/ myAvatar.png  in vscode or editor click touch pad or double click it will not zoom during scaling . 
    
# Bit Transaction 
    
    bc1q4n65rrpzz04d2ax394e0j6wmh5ayc6lvffyxc    (bitaddress)

# Codespaces
    
    https://ali2210-wisdomenigma-pixelszoom-7pr5vqxg2xrjv.github.dev/
