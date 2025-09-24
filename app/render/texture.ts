import * as THREE from 'three'

const loader = new THREE.TextureLoader()
export function load_texture(url: string) {
    return loader.load(
        url,
        (tex) => {
            console.log(`Texture ${tex} loaded successfully`)
        },
        undefined,
        (err) => {
            console.error('Error loading texture:', err)
        }
    )
}

export const alphatest_colour = (shader: THREE.WebGLProgramParametersWithUniforms) => {
    shader.fragmentShader = shader.fragmentShader.replace('#include <alphatest_fragment>', `
#ifdef USE_ALPHATEST

\t#ifdef ALPHA_TO_COVERAGE

\tdiffuseColor.a = smoothstep( alphaTest, alphaTest + fwidth( diffuseColor.a ), diffuseColor.a );
\tif ( diffuseColor.a == 0.0 ) diffuseColor.rgb = vec3( 1.0 );

\t#else

\tif ( diffuseColor.a < alphaTest ) diffuseColor.rgb = vec3( 1.0 );

\t#endif

#endif
`);
}