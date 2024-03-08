package engine

import (
	"fmt"
	"strings"

	sdf "github.com/supersdf-go/engine/sdf"
)

var (
	sdfvertexShaderSource = `
		#version 410
		uniform mat4 modelView;
		uniform mat4 model;
		
		in vec3 vp;
		out vec3 wp;
		out vec3 eye_dir;
		void main() {
			gl_Position = modelView * vec4(vp, 1.0);
			wp = (model * vec4(vp, 1.0)).xyz;
		}
	` + "\x00"

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

		float sdf(vec3 p){ 
			return // SDF_INNER ;
		}
		vec4 colorsdf(vec3 p){
			vec4 base = vec4(1,1,1,1);
			base = // SDF_COLOR_INNER ;
			return base;
		}


		void main() {
			vec3 loc = wp;
			vec3 dir = normalize(wp - cameraPosition);
			for(int i =0; i <20	 ;i++){
				float d = sdf(loc);
				loc = loc + d * 1.2 * dir;
			}

			if (sdf(loc) < 0.1) {
				frag_color = colorsdf(loc);
				
			}else{
				discard;
			}
		}
	` + "\x00"
)

func SDF2GLSL_inner(sdfObj sdf.Sdf) string {

	switch obj := sdfObj.(type) {
	case sdf.Sphere:
		return fmt.Sprintf("sphere(p,vec3(%v, %v, %v), %v)", obj.Center.X, obj.Center.Y, obj.Center.Z, obj.Radius)
	case sdf.Color:
		return SDF2GLSL_inner(obj.Sub)
	case sdf.Union:
		if len(obj) == 0 {
			return "infinity()"
		}

		str := SDF2GLSL_inner(obj[0])
		for i := 1; i < len(obj); i++ {
			str = fmt.Sprintf("min(%v, %v)", str, SDF2GLSL_inner(obj[i]))
		}
		return str
	default:
		panic(fmt.Sprintf("Unsupported type: %v", obj))
	}
}

func SDF2GLSL(sdfObj sdf.Sdf) string {
	base := sdffragmentShaderSource
	innerGlsl := SDF2GLSL_inner(sdfObj)
	result := strings.Replace(base, "// SDF_INNER", innerGlsl, 1)
	return result
}
