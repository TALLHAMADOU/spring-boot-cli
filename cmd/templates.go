package cmd

import (
	"bytes"
	"fmt"
	"text/template"
)

// renderTemplate parses and executes a named Java source template with the given data.
func renderTemplate(name string, data interface{}) (string, error) {
	tmplStr, ok := javaTemplates[name]
	if !ok {
		return "", fmt.Errorf("template %q not found", name)
	}
	tmpl, err := template.New(name).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", name, err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template %q: %w", name, err)
	}
	return buf.String(), nil
}

// javaTemplates is the registry of all Java source file templates used by generators.
var javaTemplates = map[string]string{
	"repository":                repositoryTmpl,
	"service_interface":         serviceIfaceTmpl,
	"service_interface_generic": serviceIfaceGenericTmpl,
	"service_impl":              serviceImplTmpl,
	"service_impl_generic":      serviceImplGenericTmpl,
	"controller_crud":           controllerCrudTmpl,
	"controller_basic":          controllerBasicTmpl,
	"test_service":              testServiceTmpl,
	"test_controller":           testControllerTmpl,
	"application":               applicationTmpl,
	"entity":                    entityTmpl,
}

const entityTmpl = `package {{.Pkg}}.entity;

{{range .Imports}}import {{.}};
{{end}}
{{if .Auditing}}@EntityListeners(AuditingEntityListener.class)
{{end}}@Entity
{{- if .Lombok}}
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
{{- end}}
public class {{.Name}} {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
{{- if .Auditing}}

    // Auditing fields
    @CreatedDate
    private java.time.Instant createdAt;

    @LastModifiedDate
    private java.time.Instant updatedAt;
{{- end}}
{{- range .Fields}}

    private {{.Type}} {{.Name}};
{{- end}}
{{- if not .Lombok}}

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }
{{- range .Fields}}

    public {{.Type}} get{{.Cap}}() {
        return {{.Name}};
    }

    public void set{{.Cap}}({{.Type}} {{.Name}}) {
        this.{{.Name}} = {{.Name}};
    }
{{- end}}
{{- end}}
}
`

const repositoryTmpl = `package {{.Pkg}}.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import {{.Pkg}}.entity.{{.Entity}};

public interface {{.Entity}}Repository extends JpaRepository<{{.Entity}}, Long> {
}
`

const serviceIfaceTmpl = `package {{.Pkg}}.service;

import java.util.List;
import java.util.Optional;
import {{.Pkg}}.entity.{{.Entity}};

public interface {{.Name}}Service {
    List<{{.Entity}}> findAll();
    Optional<{{.Entity}}> findById(Long id);
    {{.Entity}} save({{.Entity}} entity);
    Optional<{{.Entity}}> update(Long id, {{.Entity}} entity);
    void deleteById(Long id);
}
`

const serviceIfaceGenericTmpl = `package {{.Pkg}}.service;

import java.util.List;
import java.util.Optional;

/**
 * Generic service interface for {{.Name}}.
 * TODO: Replace Object with your actual domain type and inject the appropriate repository.
 */
public interface {{.Name}}Service {
    List<Object> findAll();
    Optional<Object> findById(Long id);
    Object save(Object entity);
    Optional<Object> update(Long id, Object entity);
    void deleteById(Long id);
}
`

const serviceImplTmpl = `package {{.Pkg}}.service.impl;

import java.util.*;
import org.springframework.stereotype.Service;
import org.springframework.beans.factory.annotation.Autowired;
import {{.Pkg}}.service.{{.Name}}Service;
import {{.Pkg}}.entity.{{.Entity}};
import {{.Pkg}}.repository.{{.Entity}}Repository;

@Service
public class {{.Name}}ServiceImpl implements {{.Name}}Service {

    private final {{.Entity}}Repository repository;

    @Autowired
    public {{.Name}}ServiceImpl({{.Entity}}Repository repository) {
        this.repository = repository;
    }

    @Override
    public List<{{.Entity}}> findAll() {
        return repository.findAll();
    }

    @Override
    public Optional<{{.Entity}}> findById(Long id) {
        return repository.findById(id);
    }

    @Override
    public {{.Entity}} save({{.Entity}} entity) {
        return repository.save(entity);
    }

    @Override
    public Optional<{{.Entity}}> update(Long id, {{.Entity}} entity) {
        return repository.findById(id).map(existing -> {
{{.CopyLines}}        });
    }

    @Override
    public void deleteById(Long id) {
        repository.deleteById(id);
    }
}
`

const serviceImplGenericTmpl = `package {{.Pkg}}.service.impl;

import java.util.*;
import org.springframework.stereotype.Service;
import {{.Pkg}}.service.{{.Name}}Service;

/**
 * Generic service implementation for {{.Name}}.
 * TODO: Inject your repository and replace Object with your actual domain type.
 */
@Service
public class {{.Name}}ServiceImpl implements {{.Name}}Service {

    @Override
    public List<Object> findAll() {
        return Collections.emptyList();
    }

    @Override
    public Optional<Object> findById(Long id) {
        return Optional.empty();
    }

    @Override
    public Object save(Object entity) {
        return entity;
    }

    @Override
    public Optional<Object> update(Long id, Object entity) {
        return Optional.empty();
    }

    @Override
    public void deleteById(Long id) {
    }
}
`

const controllerCrudTmpl = `package {{.Pkg}}.controller;

import org.springframework.web.bind.annotation.*;
import java.util.List;
import {{.Pkg}}.entity.{{.Entity}};
import {{.Pkg}}.service.{{.Entity}}Service;
import org.springframework.http.ResponseEntity;
import org.springframework.beans.factory.annotation.Autowired;

@RestController
@RequestMapping("/api/{{.EntityLower}}")
public class {{.Name}}Controller {

    @Autowired
    private {{.Entity}}Service service;

    @GetMapping
    public List<{{.Entity}}> list() {
        return service.findAll();
    }

    @GetMapping("/{id}")
    public ResponseEntity<{{.Entity}}> get(@PathVariable Long id) {
        return service.findById(id).map(ResponseEntity::ok).orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    public {{.Entity}} create(@RequestBody {{.Entity}} entity) {
        return service.save(entity);
    }

    @PutMapping("/{id}")
    public ResponseEntity<{{.Entity}}> update(@PathVariable Long id, @RequestBody {{.Entity}} entity) {
        return service.update(id, entity).map(ResponseEntity::ok).orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable Long id) {
        service.deleteById(id);
        return ResponseEntity.noContent().build();
    }
}
`

const controllerBasicTmpl = `package {{.Pkg}}.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/{{.NameLower}}")
public class {{.Name}}Controller {

    @GetMapping
    public String index() {
        return "ok";
    }
}
`

const testServiceTmpl = `package {{.Pkg}}.service;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import static org.assertj.core.api.Assertions.assertThat;
import {{.Pkg}}.service.impl.{{.Name}}ServiceImpl;

@ExtendWith(MockitoExtension.class)
public class {{.Name}}ServiceTest {

    @Mock
    // TODO: Replace Object with your actual repository type, e.g.: {{.Name}}Repository repository;
    private Object repository;

    @InjectMocks
    private {{.Name}}ServiceImpl service;

    @Test
    void testFindAll() {
        // TODO: implement test
        // assertThat(service.findAll()).isEmpty();
    }
}
`

const testControllerTmpl = `package {{.Pkg}}.controller;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.beans.factory.annotation.Autowired;
import static org.assertj.core.api.Assertions.assertThat;

@WebMvcTest({{.Name}}Controller.class)
public class {{.Name}}ControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    // TODO: Replace Object with your actual service type, e.g.: {{.Name}}Service service;
    private Object service;

    @Test
    void testGetAll() throws Exception {
        // TODO: implement MockMvc test
        assertThat(mockMvc).isNotNull();
    }
}
`

const applicationTmpl = `package {{.Pkg}};

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
`
