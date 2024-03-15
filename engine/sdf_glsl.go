package engine

import (
	"fmt"
	"strings"

	sdf "github.com/supersdf-go/engine/sdf"
)

var (
	sdffragmentShaderSource = `
		#version 410

		uniform vec4 color;
		out vec4 frag_color;

		uniform vec3 cameraPosition;
		in vec3 wp;
		in vec3 eye_dir;

		float sphere(vec3 p, vec3 c, float r){
			return length(p - c) -r;
		}

		void sdf(vec3 p, inout float outdist, inout vec4 outcolor){ 
			vec4 color = vec4(1,1,1,1);
			float d = 100000.0;
			// SDF_INNER
			outdist = d;
			outcolor = color;
		}

		void main() {
			vec3 loc = wp;
			vec3 dir = normalize(wp - cameraPosition);
			float adist;
			vec4 acolor;
			for(int i =0; i <20	 ;i++){
				sdf(loc,adist,acolor);
				
				loc = loc + adist * 1.2 * dir;
			}

			if (adist < 0.1) {
				frag_color = vec4(1.0,0.1,0.1,1);
				//frag_color = acolor;
				
			}else{
				frag_color = vec4(0.1,0.1,0.1,1);
				//discard;
			}
		}
	` + "\x00"
)

func SDF2GLSL_inner(sdfObj sdf.Sdf, output *string) {

	switch obj := sdfObj.(type) {
	case sdf.Sphere:

		*output = fmt.Sprintf("%v\nd = sphere(p,vec3(%v, %v, %v), %v);", *output, obj.Center.X, obj.Center.Y, obj.Center.Z, obj.Radius)
	case sdf.Color:
		*output = fmt.Sprintf("%v\ncolor = vec4(%v, %v, %v, 1);", *output, obj.Color.X, obj.Color.Y, obj.Color.Z)
		SDF2GLSL_inner(obj.Sub, output)
	case sdf.Union:
		if len(obj) == 0 {
			return
		}
		if len(obj) == 1 {
			SDF2GLSL_inner(obj[0], output)
		}
		SDF2GLSL_inner(obj[0], output)

		for i := 1; i < len(obj); i++ {
			inner := ""
			SDF2GLSL_inner(obj[i], &inner)
			*output = fmt.Sprintf("%v{float d2 = d;vec4 color2 = color;  %v if(d > d2){d = d2; color = color2;}}", *output, inner)

		}
	default:
		panic(fmt.Sprintf("Unsupported type: %v", obj))
	}
}

func SDF2GLSL(sdfObj sdf.Sdf) string {
	base := sdffragmentShaderSource
	result := ""
	SDF2GLSL_inner(sdfObj, &result)
	result2 := strings.Replace(base, "// SDF_INNER", result, 1)
	result2 = strings.Replace(result2, "// SDF_COLOR_INNER", result, 1)
	return result2
}
