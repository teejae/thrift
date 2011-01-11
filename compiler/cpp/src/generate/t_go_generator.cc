/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

#include <string>
#include <fstream>
#include <iostream>
#include <vector>
#include <stdlib.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <sstream>
#include <algorithm>
#include "t_generator.h"
#include "platform.h"
using namespace std;

/**
 * Go code generator.
 *
 */
class t_go_generator : public t_generator {
 public:
  t_go_generator(
      t_program* program,
      const std::map<std::string, std::string>& parsed_options,
      const std::string& option_string)
    : t_generator(program)
  {
		out_dir_base_ = "gen-go";
    std::map<std::string, std::string>::const_iterator iter;

    iter = parsed_options.find("spaces");
    if (iter != parsed_options.end()) {
      indent_str_ = "  ";
    } else {
      indent_str_ = "\t";
    }
  }

  /**
   * Init and close methods
   */

  void init_generator();
  void close_generator();

  /**
   * Program-level generation functions
   */

  void generate_typedef  (t_typedef*  ttypedef);
  void generate_enum     (t_enum*     tenum);
  void generate_const    (t_const*    tconst);
  void generate_struct   (t_struct*   tstruct);
  void generate_xception (t_struct*   txception);
  void generate_service  (t_service*  tservice);

  std::string render_const_value(t_type* type, t_const_value* value);

  /**
   * Struct generation code
   */

  void generate_go_struct(t_struct* tstruct, bool is_exception);
  void generate_go_struct_definition(std::ofstream& out, t_struct* tstruct, bool is_xception=false, bool is_result=false);
  void generate_go_struct_constructor(std::ofstream& out, t_struct* tstruct);
  void generate_go_struct_reader(std::ofstream& out, t_struct* tstruct);
  void generate_go_struct_writer(std::ofstream& out, t_struct* tstruct);
  void generate_go_struct_required_validator(std::ofstream& out, t_struct* tstruct);
  void generate_py_function_helpers(t_function* tfunction);

  /**
   * Service-level generation functions
   */

  void generate_service_helpers   (t_service*  tservice);
  void generate_service_interface (t_service* tservice);
  void generate_service_client    (t_service* tservice);
  void generate_service_remote    (t_service* tservice);
  void generate_service_server    (t_service* tservice);
  void generate_process_function  (t_service* tservice, t_function* tfunction);

  /**
   * Serialization constructs
   */

  void generate_deserialize_field        (std::ofstream &out,
                                          t_field*    tfield,
                                          std::string prefix="",
                                          bool inclass=false);

  void generate_deserialize_struct       (std::ofstream &out,
                                          t_struct*   tstruct,
                                          std::string prefix="");

  void generate_deserialize_container    (std::ofstream &out,
                                          t_type*     ttype,
                                          std::string prefix="");

  void generate_deserialize_set_element  (std::ofstream &out,
                                          t_set*      tset,
                                          std::string prefix="");

  void generate_deserialize_map_element  (std::ofstream &out,
                                          t_map*      tmap,
                                          std::string prefix="");

  void generate_deserialize_list_element (std::ofstream &out,
                                          t_list*     tlist,
                                          std::string prefix="");

  void generate_serialize_field          (std::ofstream &out,
                                          t_field*    tfield,
                                          std::string prefix="");

  void generate_serialize_struct         (std::ofstream &out,
                                          t_struct*   tstruct,
                                          std::string prefix="");

  void generate_serialize_container      (std::ofstream &out,
                                          t_type*     ttype,
                                          std::string prefix="");

  void generate_serialize_map_element    (std::ofstream &out,
                                          t_map*      tmap,
                                          std::string kiter,
                                          std::string viter);

  void generate_serialize_set_element    (std::ofstream &out,
                                          t_set*      tmap,
                                          std::string iter);

  void generate_serialize_list_element   (std::ofstream &out,
                                          t_list*     tlist,
                                          std::string iter);

  void generate_python_docstring         (std::ofstream& out,
                                          t_struct* tstruct);

  void generate_python_docstring         (std::ofstream& out,
                                          t_function* tfunction);

  void generate_python_docstring         (std::ofstream& out,
                                          t_doc*    tdoc,
                                          t_struct* tstruct,
                                          const char* subheader);

  void generate_python_docstring         (std::ofstream& out,
                                          t_doc* tdoc);

  /**
   * Helper rendering functions
   */

  std::string go_autogen_comment();
  std::string go_package();
  std::string go_imports();
  std::string render_includes();
  std::string declare_argument(t_field* tfield);
  std::string render_field_default_value(t_field* tfield);
  std::string type_name(t_type* ttype);
  std::string function_signature(t_function* tfunction);
  std::string function_signature_if(t_function* tfunction);
  std::string argument_list(t_struct* tstruct);
  std::string type_to_enum(t_type* ttype);
  std::string type_to_spec_args(t_type* ttype);

  std::string base_type_name(t_base_type*);

  static bool is_valid_namespace(const std::string& sub_namespace) {
    return sub_namespace == "twisted";
  }

  static std::string get_real_py_module(const t_program* program, bool gen_twisted) {
    if(gen_twisted) {
      std::string twisted_module = program->get_namespace("py.twisted");
      if(!twisted_module.empty()){
        return twisted_module;
      }
    }

    std::string real_module = program->get_namespace("py");
    if (real_module.empty()) {
      return program->get_name();
    }
    return real_module;
  }

 private:
  /**
   * True iff we should generate new-style classes.
   */
  bool gen_newstyle_;

  /**
   * True iff we should generate Twisted-friendly RPC services.
   */
  bool gen_twisted_;

  /**
   * True iff strings should be encoded using utf-8.
   */
  bool gen_utf8strings_;

  /**
   * File streams
   */

  std::ofstream f_types_;
  std::ofstream f_consts_;
  std::ofstream f_service_;

  std::string package_dir_;

};


/**
 * Prepares for file generation by opening up the necessary file output
 * streams.
 *
 * @param tprogram The program to generate
 */
void t_go_generator::init_generator() {
  // Make output directory
  string package_name_ = program_->get_namespace("as3");

  string module = get_real_py_module(program_, gen_twisted_);
  MKDIR(get_out_dir().c_str());
  package_dir_ = get_out_dir() + "/" + program_->get_name();
	MKDIR(package_dir_.c_str());

  // Make output file
  string f_types_name = package_dir_+"/"+"ttypes.go";
  f_types_.open(f_types_name.c_str());

  string f_consts_name = package_dir_+"/"+"constants.go";
  f_consts_.open(f_consts_name.c_str());

  // Print header
  f_types_ <<
    go_autogen_comment() << endl <<
    go_package() << endl <<
    go_imports() << endl <<
    render_includes() << endl;

  f_consts_ <<
    go_autogen_comment() << endl <<
    go_package() << endl <<
		"const (" << endl;
}

/**
 * Renders all the imports necessary for including another Thrift program
 */
string t_go_generator::render_includes() {
  const vector<t_program*>& includes = program_->get_includes();

  string result = "";
  if (includes.size() == 0) {
    result += "\n";
  } else {
    result += "import (\n";
    indent_up();
    for (size_t i = 0; i < includes.size(); ++i) {
      result +=  indent() + string("\"") + get_real_py_module(includes[i], gen_twisted_) + "\"\n";
    }
    indent_down();
    result += ")\n";
  }
  return result;
}

/**
 * Autogen'd comment
 */
string t_go_generator::go_autogen_comment() {
  return
    std::string("//\n") +
    "// Autogenerated by Thrift\n" +
    "//\n" +
    "// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING\n" +
    "//\n";
}

/**
 * Prints standard thrift imports
 */
string t_go_generator::go_package() {
  return string("package ") + program_->get_name() + "\n";
}

/**
 * Prints standard thrift imports
 */
string t_go_generator::go_imports() {
  return
		string("import \"github.com/teejae/go-thrift/thrift\"");
}

/**
 * Closes the type files
 */
void t_go_generator::close_generator() {
  // Close types file
  f_types_.close();
	f_consts_ << ")" << endl;
  f_consts_.close();
}

/**
 * Generates a typedef.
 *
 * @param ttypedef The type definition
 */
void t_go_generator::generate_typedef(t_typedef* ttypedef) {
  f_types_ <<
    indent() << "type " << ttypedef->get_symbolic () << " " << type_name(ttypedef->get_type()) << endl;
}

/**
 * Generates code for an enumerated type. Done as a using a class to scope
 * the values.
 *
 *
 * 
enum Operation {
  ADD = 1,
  SUBTRACT = 2,
  MULTIPLY = 3,
  DIVIDE = 4
}

 *	type Operation int32

 * const (
		OPERATION_ADD      Operation = 1
		OPERATION_SUBTRACT Operation = 2
		OPERATION_MULTIPLY Operation = 3
		OPERATION_DIVIDE   Operation = 4
 * )

 * @param tenum The enumeration
 */
void t_go_generator::generate_enum(t_enum* tenum) {
 	f_types_ <<
		"type " << tenum->get_name() <<
		" int32" << endl <<
		"const (" << endl;
		
  indent_up();
 
  vector<t_enum_value*> constants = tenum->get_constants();
  vector<t_enum_value*>::iterator c_iter;
  for (c_iter = constants.begin(); c_iter != constants.end(); ++c_iter) {
    int value = (*c_iter)->get_value();
		// FIXME: make this to_upper_case, once figure out how to call it
		string const_name = capitalize(tenum->get_name()) + "_" + (*c_iter)->get_name();
    indent(f_types_) << const_name << " " << tenum->get_name() << " = " << value << endl;
  }
 
  indent_down();
	f_types_ << ")" << endl;
}

/**
 * Generate a constant value
 */
void t_go_generator::generate_const(t_const* tconst) {
  t_type* type = tconst->get_type();
  string name = tconst->get_name();
  t_const_value* value = tconst->get_value();

  indent(f_consts_) << name << " = " << render_const_value(type, value);
  f_consts_ << endl;
}

/**
 * Prints the value of a constant with the given type. Note that type checking
 * is NOT performed in this function as it is always run beforehand using the
 * validate_types method in main.cc
 */
string t_go_generator::render_const_value(t_type* type, t_const_value* value) {
  type = get_true_type(type);
  std::ostringstream out;

  if (type->is_base_type()) {
    t_base_type::t_base tbase = ((t_base_type*)type)->get_base();
    switch (tbase) {
    case t_base_type::TYPE_STRING:
      out << '"' << get_escaped_string(value) << '"';
      break;
    case t_base_type::TYPE_BOOL:
      out << (value->get_integer() > 0 ? "true" : "false");
      break;
    case t_base_type::TYPE_BYTE:
			out << "byte(" << value->get_integer() << ")";
      break;
    case t_base_type::TYPE_I16:
			out << "int16(" << value->get_integer() << ")";
      break;
    case t_base_type::TYPE_I32:
			out << "int32(" << value->get_integer() << ")";
      break;
    case t_base_type::TYPE_I64:
			out << "int64(" << value->get_integer() << ")";
      break;
    case t_base_type::TYPE_DOUBLE:
    	out << "float64(" << value->get_integer() << ")";
      break;
    default:
      throw "compiler error: no const of base type " + t_base_type::t_base_name(tbase);
    }
  } else if (type->is_enum()) {
    indent(out) << value->get_integer();
  } else if (type->is_struct() || type->is_xception()) {
		// FIXME: finish this
    out << type->get_name() << "(**{" << endl;
    indent_up();
    const vector<t_field*>& fields = ((t_struct*)type)->get_members();
    vector<t_field*>::const_iterator f_iter;
    const map<t_const_value*, t_const_value*>& val = value->get_map();
    map<t_const_value*, t_const_value*>::const_iterator v_iter;
    for (v_iter = val.begin(); v_iter != val.end(); ++v_iter) {
      t_type* field_type = NULL;
      for (f_iter = fields.begin(); f_iter != fields.end(); ++f_iter) {
        if ((*f_iter)->get_name() == v_iter->first->get_string()) {
          field_type = (*f_iter)->get_type();
        }
      }
      if (field_type == NULL) {
        throw "type error: " + type->get_name() + " has no field " + v_iter->first->get_string();
      }
      out << indent();
      out << render_const_value(g_type_string, v_iter->first);
      out << " : ";
      out << render_const_value(field_type, v_iter->second);
      out << "," << endl;
    }
    indent_down();
    indent(out) << "})";
  } else if (type->is_map()) {
    t_type* ktype = ((t_map*)type)->get_key_type();
    t_type* vtype = ((t_map*)type)->get_val_type();
    out << "map[" << ktype->get_name() << "]" << vtype->get_name() << " {" << endl;
    indent_up();
    const map<t_const_value*, t_const_value*>& val = value->get_map();
    map<t_const_value*, t_const_value*>::const_iterator v_iter;
    for (v_iter = val.begin(); v_iter != val.end(); ++v_iter) {
      out << indent();
      out << render_const_value(ktype, v_iter->first);
      out << " : ";
      out << render_const_value(vtype, v_iter->second);
      out << "," << endl;
    }
    indent_down();
    indent(out) << "}";
  } else if (type->is_list() || type->is_set()) {
    t_type* etype;
    if (type->is_list()) {
      etype = ((t_list*)type)->get_elem_type();
    } else {
      etype = ((t_set*)type)->get_elem_type();
    }
    if (type->is_set()) {
      out << "set(";
    }
    out << "[" << endl;
    indent_up();
    const vector<t_const_value*>& val = value->get_list();
    vector<t_const_value*>::const_iterator v_iter;
    for (v_iter = val.begin(); v_iter != val.end(); ++v_iter) {
      out << indent();
      out << render_const_value(etype, *v_iter);
      out << "," << endl;
    }
    indent_down();
    indent(out) << "]";
    if (type->is_set()) {
      out << ")";
    }
  } else {
    throw "CANNOT GENERATE CONSTANT FOR TYPE: " + type->get_name();
  }

  return out.str();
}

/**
 * Generates a python struct
 */
void t_go_generator::generate_struct(t_struct* tstruct) {
  generate_go_struct(tstruct, false);
}

/**
 * Generates a struct definition for a thrift exception. Basically the same
 * as a struct but extends the Exception class.
 *
 * @param txception The struct definition
 */
void t_go_generator::generate_xception(t_struct* txception) {
  generate_go_struct(txception, true);
}

/**
 * Generates a go struct
 */
void t_go_generator::generate_go_struct(t_struct* tstruct,
                                        bool is_exception) {
  generate_go_struct_definition(f_types_, tstruct, is_exception);
}

/**
 * Generates a struct definition for a thrift data type.

type SharedStruct struct {
	Key   *int32
	Value *string
}
 *
 * @param tstruct The struct definition
 */
void t_go_generator::generate_go_struct_definition(ofstream& out,
                                                   t_struct* tstruct,
                                                   bool is_exception,
                                                   bool is_result) {

  const vector<t_field*>& members = tstruct->get_members();
  vector<t_field*>::const_iterator m_iter;

  out << std::endl;
  generate_python_docstring(out, tstruct);
  out << "type " << capitalize(tstruct->get_name()) << " struct {" << endl;

  indent_up();

  if (is_exception) {
    indent(out) << "*thrift.TException" << endl;
  }

  if (members.size() > 0) {
    for (m_iter = members.begin(); m_iter != members.end(); ++m_iter) {
      indent(out) << capitalize((*m_iter)->get_name()) + " *" + type_name((*m_iter)->get_type()) << endl;
    }
  }

  indent_down();
  out << "}" << endl;

  generate_go_struct_constructor(out, tstruct);
  generate_go_struct_reader(out, tstruct);
  generate_go_struct_writer(out, tstruct);
}

/**
 * Generates the constructor for a struct
 * func NewSharedStruct() *SharedStruct {
 * 	return &SharedStruct{}
 * }
 */
void t_go_generator::generate_go_struct_constructor(ofstream& out,
                                                    t_struct* tstruct) {
  string name = capitalize(tstruct->get_name());
  indent(out) << "func New" << name << "() *" << name << " {" << endl;
  indent_up();
  indent(out) << "return &" << name << "{}" << endl;
  indent_down();
  indent(out) << "}" << endl;
}

/**
 * Generates the read method for a struct
 
 func (s *SharedStruct) Read(iprot thrift.TProtocol) {
 	iprot.ReadStructBegin()
 	for {
 		_, ftype, fid := iprot.ReadFieldBegin()
 		if ftype == thrift.TTYPE_STOP {
 			break
 		}

 		switch fid {
 		case 1: // key
 			if ftype == thrift.TTYPE_I32 {
 				v := iprot.ReadI32()
 				s.Key = &v
 			} else {
 				thrift.SkipType(iprot, ftype)
 			}
 		case 2: // value
 			if ftype == thrift.TTYPE_STRING {
 				v := iprot.ReadString()
 				s.Value = &v
 			} else {
 				thrift.SkipType(iprot, ftype)
 			}
 		default: // unknown
 			thrift.SkipType(iprot, ftype)
 		}
 		iprot.ReadFieldEnd()
 	}
 	iprot.ReadStructEnd()
 }
 
 */
void t_go_generator::generate_go_struct_reader(ofstream& out,
                                                t_struct* tstruct) {
  const vector<t_field*>& fields = tstruct->get_members();
  vector<t_field*>::const_iterator f_iter;

  indent(out) <<
    "func (s *" << capitalize(tstruct->get_name()) << ") Read(iprot thrift.TProtocol) {" << endl;
  indent_up();

  indent(out) <<
    "iprot.ReadStructBegin()" << endl;

  // Loop over reading in fields
  indent(out) <<
    "for {" << endl;
  indent_up();

  // Read beginning field marker
  indent(out) <<
    "_, ftype, fid := iprot.ReadFieldBegin()" << endl;

  // Check for field STOP marker and break
  indent(out) <<
    "if ftype == thrift.TTYPE_STOP {" << endl;
  indent_up();
  indent(out) <<
    "break" << endl;
  indent_down();
  indent(out) <<
    "}" << endl << endl;

  // Switch statement on the field we are reading
  indent(out) <<
    "switch fid {" << endl;

  // Generate deserialization code for known cases
  for (f_iter = fields.begin(); f_iter != fields.end(); ++f_iter) {
    indent(out) << "case " << (*f_iter)->get_key() << ":" << endl;
    indent_up();
    indent(out) << "if ftype == " << type_to_enum((*f_iter)->get_type()) << " {" << endl;
    indent_up();
    generate_deserialize_field(out, *f_iter, "s.");
    indent_down();
    indent(out) << "} else {" << endl;
    indent_up();
    out <<
      indent() << "thrift.SkipType(iprot, ftype)" << endl;
    indent_down();
    out <<
      indent() << "}" << endl;
    indent_down();
  }

  // In the default case we skip the field
  out <<
    indent() << "default:" << endl;
  indent_up();
  out <<
    indent() << "thrift.SkipType(iprot, ftype)" << endl;
  indent_down();

  indent(out) << "}" << endl;

  // Read field end marker
  indent(out) <<
    "iprot.ReadFieldEnd()" << endl;

  indent_down();

  indent(out) << "}" << endl;
  indent(out) << "iprot.ReadStructEnd()" << endl;

  indent_down();
  
  out << "}" << endl;
}

void t_go_generator::generate_go_struct_writer(ofstream& out,
                                               t_struct* tstruct) {
  string name = tstruct->get_name();
  const vector<t_field*>& fields = tstruct->get_sorted_members();
  vector<t_field*>::const_iterator f_iter;

  indent(out) <<
    "func (s *" << capitalize(tstruct->get_name()) << ") Write(oprot thrift.TProtocol) {" << endl;
  indent_up();

  indent(out) <<
    "oprot.WriteStructBegin(\"" << name << "\")" << endl;

  for (f_iter = fields.begin(); f_iter != fields.end(); ++f_iter) {
    // Write field header
    indent(out) <<
      "if s." << capitalize((*f_iter)->get_name()) << " != nil {" << endl;
    indent_up();
    indent(out) <<
      "oprot.WriteFieldBegin(" <<
      "\"" << (*f_iter)->get_name() << "\", " <<
      type_to_enum((*f_iter)->get_type()) << ", " <<
      (*f_iter)->get_key() << ")" << endl;

    // Write field contents
    generate_serialize_field(out, *f_iter, "s.");

    // Write field closer
    indent(out) <<
      "oprot.WriteFieldEnd()" << endl;

    indent_down();
    indent(out) << "}" << endl;
  }

  // Write the struct map
  out <<
    indent() << "oprot.WriteFieldStop()" << endl <<
    indent() << "oprot.WriteStructEnd()" << endl;

  // TODO: do we need validator?
  // generate_go_struct_required_validator(out, tstruct);

  indent_down();
  out << "}" << endl;
}

void t_go_generator::generate_go_struct_required_validator(ofstream& out,
                                               t_struct* tstruct) {
  indent(out) << "def validate(self):" << endl;
  indent_up();

  const vector<t_field*>& fields = tstruct->get_members();

  if (fields.size() > 0) {
    vector<t_field*>::const_iterator f_iter;

    for (f_iter = fields.begin(); f_iter != fields.end(); ++f_iter) {
      t_field* field = (*f_iter);
      if (field->get_req() == t_field::T_REQUIRED) {
        indent(out) << "if self." << field->get_name() << " is None:" << endl;
        indent(out) << "  raise TProtocol.TProtocolException(message='Required field " <<
          field->get_name() << " is unset!')" << endl;
      }
    }
  }

  indent(out) << "return" << endl << endl;
  indent_down();
}

/**
 * Generates a thrift service.
 *
 * @param tservice The service definition
 */
void t_go_generator::generate_service(t_service* tservice) {
  string f_service_name = package_dir_+ "/" + lowercase(service_name_) + ".go";
  f_service_.open(f_service_name.c_str());

  f_service_ <<
    go_autogen_comment() << endl <<
    go_package() << endl <<
    go_imports() << endl;

  if (tservice->get_extends() != NULL) {
    f_service_ <<
      "import \"" << tservice->get_extends()->get_program()->get_name() << "\"" << endl;
  }

  f_service_ <<
    "import \"fmt\"";

  f_service_ << endl;

  // Generate the three main parts of the service (well, two for now in PHP)
  generate_service_interface(tservice);
  generate_service_client(tservice);
  generate_service_server(tservice);
  generate_service_helpers(tservice);
  generate_service_remote(tservice);

  // Close service file
  f_service_.close();
}

/**
 * Generates helper functions for a service.
 *
 * @param tservice The service to generate a header definition for
 */
void t_go_generator::generate_service_helpers(t_service* tservice) {
  vector<t_function*> functions = tservice->get_functions();
  vector<t_function*>::iterator f_iter;

  f_service_ <<
    "// HELPER FUNCTIONS AND STRUCTURES" << endl;

  for (f_iter = functions.begin(); f_iter != functions.end(); ++f_iter) {
    t_struct* ts = (*f_iter)->get_arglist();
    generate_go_struct_definition(f_service_, ts, false);
    generate_py_function_helpers(*f_iter);
  }
}

/**
 * Generates a struct and helpers for a function.
 *
 * @param tfunction The function
 */
void t_go_generator::generate_py_function_helpers(t_function* tfunction) {
  if (!tfunction->is_oneway()) {
    t_struct result(program_, tfunction->get_name() + "_result");
    t_field success(tfunction->get_returntype(), "success", 0);
    if (!tfunction->get_returntype()->is_void()) {
      result.append(&success);
    }

    t_struct* xs = tfunction->get_xceptions();
    const vector<t_field*>& fields = xs->get_members();
    vector<t_field*>::const_iterator f_iter;
    for (f_iter = fields.begin(); f_iter != fields.end(); ++f_iter) {
      result.append(*f_iter);
    }
    generate_go_struct_definition(f_service_, &result, false, true);
  }
}

/**
 * Generates a service interface definition.
 *
 type SharedService interface {
 	GetStruct(key int32) *SharedStruct
 }
 
 * @param tservice The service to generate a header definition for
 */
void t_go_generator::generate_service_interface(t_service* tservice) {
  string extends = "";
  if (tservice->get_extends() != NULL) {
    extends = tservice->get_extends()->get_program()->get_name() + "." + type_name(tservice->get_extends());
  }

  // FIXME: deal with extends
  generate_python_docstring(f_service_, tservice);
  f_service_ <<
    "type " << tservice->get_name() << " interface {" << endl;
  indent_up();
  indent(f_service_) << extends << endl;
  vector<t_function*> functions = tservice->get_functions();
  vector<t_function*>::iterator f_iter;
  for (f_iter = functions.begin(); f_iter != functions.end(); ++f_iter) {
    generate_python_docstring(f_service_, (*f_iter));
    f_service_ <<
      indent() << function_signature_if(*f_iter) << endl;
    indent_up();
    indent_down();
  }

  indent_down();
  f_service_ << "}" << endl;
}

/**
 * Generates a service client definition.
 type SharedServiceClient struct {
 	iprot thrift.TProtocol
 	oprot thrift.TProtocol
 	seqid int32
 }
 
 *
 * @param tservice The service to generate a server for.
 */
void t_go_generator::generate_service_client(t_service* tservice) {
  string extends = "";
  string extends_client = "";
  // FIXME: deal with extends
  if (tservice->get_extends() != NULL) {
    extends = type_name(tservice->get_extends());
    if (gen_twisted_) {
      extends_client = "(" + extends + ".Client)";
    } else {
      extends_client = extends + ".Client, ";
    }
  } else {
    if (gen_twisted_ && gen_newstyle_) {
        extends_client = "(object)";
    }
  }

  generate_python_docstring(f_service_, tservice);
  
  string client_name = tservice->get_name() + "Client";
  f_service_ <<
    "type " << client_name << " struct {" << endl;

  indent_up();
  
  f_service_ <<
    indent() << "iprot thrift.TProtocol" << endl <<
    indent() << "oprot thrift.TProtocol" << endl <<
    indent() << "seqid int32" << endl;
  indent_down();

  indent(f_service_) << "}" <<  endl;

  // Constructor function
  f_service_ << "func New" << client_name << "(iprot, oprot thrift.TProtocol) *" << client_name << " {" << endl;
  indent_up();
  indent(f_service_) << "return &" << client_name << "{iprot: iprot, oprot: oprot}" << endl;
  indent_down();
  f_service_ << "}" << endl;
    

  // Generate client method implementations
  vector<t_function*> functions = tservice->get_functions();
  vector<t_function*>::const_iterator f_iter;
  for (f_iter = functions.begin(); f_iter != functions.end(); ++f_iter) {
    t_struct* arg_struct = (*f_iter)->get_arglist();
    const vector<t_field*>& fields = arg_struct->get_members();
    vector<t_field*>::const_iterator fld_iter;
    string funname = capitalize((*f_iter)->get_name());

    // Open function
    generate_python_docstring(f_service_, (*f_iter));
    indent(f_service_) << "func (s *" << client_name << ") " <<
      function_signature(*f_iter) <<
      " {" << endl;
    indent_up();
    indent(f_service_) << "s.seqid++" << endl;

    indent(f_service_) <<
      "s.Send_" << funname << "(";

    bool first = true;
    for (fld_iter = fields.begin(); fld_iter != fields.end(); ++fld_iter) {
      if (first) {
        first = false;
      } else {
        f_service_ << ", ";
      }
      f_service_ << (*fld_iter)->get_name();
    }
    f_service_ << ")" << endl;

    if (!(*f_iter)->is_oneway()) {
      f_service_ << indent();
      if (!(*f_iter)->get_returntype()->is_void()) {
        f_service_ << "return ";
      }
      f_service_ <<
        "s.Recv_" << funname << "()" << endl;
    } else {
      // nothing to do
    }
    indent_down();
    f_service_ << "}" << endl;

    indent(f_service_) <<
      "func (s *" << client_name << ") Send_" << capitalize((*f_iter)->get_name()) << "(" << argument_list((*f_iter)->get_arglist()) << ") {" << endl;

    indent_up();

    std::string argsname = capitalize((*f_iter)->get_name()) + "_args";

    // Serialize the request header
    f_service_ <<
      indent() << "s.oprot.WriteMessageBegin(\"" << (*f_iter)->get_name() << "\", thrift.TMESSAGETYPE_CALL, s.seqid)" << endl;
      
    f_service_ <<
      indent() << "args := New" << argsname << "()" << endl;

    for (fld_iter = fields.begin(); fld_iter != fields.end(); ++fld_iter) {
      f_service_ <<
        indent() << "args." << capitalize((*fld_iter)->get_name()) << " = &" << (*fld_iter)->get_name() << endl;
    }

    // Write to the stream
    f_service_ <<
      indent() << "args.Write(s.oprot)" << endl <<
      indent() << "s.oprot.WriteMessageEnd()" << endl <<
      indent() << "s.oprot.GetTransport().Flush()" << endl;

    indent_down();
    f_service_ << "}" << endl;

    // receive function
    if (!(*f_iter)->is_oneway()) {
      std::string resultname = capitalize((*f_iter)->get_name()) + "_result";
      // Open function
      f_service_ <<
        endl;
      t_struct noargs(program_);
      t_function recv_function((*f_iter)->get_returntype(),
                             string("Recv_") + capitalize((*f_iter)->get_name()),
                             &noargs);
      f_service_ <<
        indent() << "func (s *" << client_name << ") " << function_signature(&recv_function) << " {" << endl;
      indent_up();

      // TODO(mcslee): Validate message reply here, seq ids etc.

      f_service_ <<
        indent() << "_, mtype, _ := s.iprot.ReadMessageBegin()" << endl <<
        indent() << "defer s.iprot.ReadMessageEnd()" << endl;

      f_service_ <<
        indent() << "if mtype == thrift.TMESSAGETYPE_EXCEPTION {" << endl;

      indent_up();
      f_service_ <<
        indent() << "x := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_UNKNOWN, \"\")" << endl <<
        indent() << "x.Read(s.iprot)" << endl <<
        indent() << "return nil, x" << endl;
      indent_down();
      
      f_service_ <<
        indent() << "}" << endl <<
        indent() << "result := New" << resultname << "()" << endl <<
        indent() << "result.Read(s.iprot)" << endl;
      
      // Careful, only return _result if not a void function
      if (!(*f_iter)->get_returntype()->is_void()) {
        f_service_ <<
          indent() << "if result.Success != nil {" << endl;
        indent_up();
        indent(f_service_) << "return result.Success, nil" << endl;
        indent_down();
        indent(f_service_) << "}" << endl;
      }

      // FIXME: deal with exceptions
      t_struct* xs = (*f_iter)->get_xceptions();
      const std::vector<t_field*>& xceptions = xs->get_members();
      vector<t_field*>::const_iterator x_iter;
      for (x_iter = xceptions.begin(); x_iter != xceptions.end(); ++x_iter) {
        string exception_name = capitalize((*x_iter)->get_name());
        f_service_ <<
          indent() << "if result." << exception_name << " != nil {" << endl;
        indent_up();
        indent(f_service_) << "return nil, result." << exception_name << endl;
        indent_down();
        indent(f_service_) << "}" << endl;
      }

      // Careful, only return _result if not a void function
      if ((*f_iter)->get_returntype()->is_void()) {
        indent(f_service_) <<
          "return" << endl;
      } else {
        f_service_ <<
          indent() << "x := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_MISSING_RESULT, \"" <<
          (*f_iter)->get_name() << " failed: unknown result\")" << endl <<
          indent() << "return nil, x" << endl;
      }

      // Close function
      indent_down();
      f_service_ << "}" << endl;
    }
  }

  f_service_ <<
    endl;
}

/**
 * Generates a command line tool for making remote requests
 *
 * @param tservice The service to generate a remote for.
 */
void t_go_generator::generate_service_remote(t_service* tservice) {
  vector<t_function*> functions = tservice->get_functions();
  vector<t_function*>::iterator f_iter;

  string f_remote_name = package_dir_+"/"+service_name_+"-remote";
  ofstream f_remote;
  f_remote.open(f_remote_name.c_str());

  f_remote <<
    go_autogen_comment() << endl <<
    go_package() << endl <<
    "import sys" << endl <<
    "import pprint" << endl <<
    "from urlparse import urlparse" << endl <<
    "from thrift.transport import TTransport" << endl <<
    "from thrift.transport import TSocket" << endl <<
    "from thrift.transport import THttpClient" << endl <<
    "from thrift.protocol import TBinaryProtocol" << endl <<
    endl;

  f_remote <<
    "import " << service_name_ << endl <<
    "from ttypes import *" << endl <<
    endl;

  f_remote <<
    "if len(sys.argv) <= 1 or sys.argv[1] == '--help':" << endl <<
    "  print ''" << endl <<
    "  print 'Usage: ' + sys.argv[0] + ' [-h host:port] [-u url] [-f[ramed]] function [arg1 [arg2...]]'" << endl <<
    "  print ''" << endl <<
    "  print 'Functions:'" << endl;
  for (f_iter = functions.begin(); f_iter != functions.end(); ++f_iter) {
    f_remote <<
      "  print '  " << (*f_iter)->get_returntype()->get_name() << " " << (*f_iter)->get_name() << "(";
    t_struct* arg_struct = (*f_iter)->get_arglist();
    const std::vector<t_field*>& args = arg_struct->get_members();
    vector<t_field*>::const_iterator a_iter;
    int num_args = args.size();
    bool first = true;
    for (int i = 0; i < num_args; ++i) {
      if (first) {
        first = false;
      } else {
        f_remote << ", ";
      }
      f_remote <<
        args[i]->get_type()->get_name() << " " << args[i]->get_name();
    }
    f_remote << ")'" << endl;
  }
  f_remote <<
    "  print ''" << endl <<
    "  sys.exit(0)" << endl <<
    endl;

  f_remote <<
    "pp = pprint.PrettyPrinter(indent = 2)" << endl <<
    "host = 'localhost'" << endl <<
    "port = 9090" << endl <<
    "uri = ''" << endl <<
    "framed = False" << endl <<
    "http = False" << endl <<
    "argi = 1" << endl <<
    endl <<
    "if sys.argv[argi] == '-h':" << endl <<
    "  parts = sys.argv[argi+1].split(':')" << endl <<
    "  host = parts[0]" << endl <<
    "  port = int(parts[1])" << endl <<
    "  argi += 2" << endl <<
    endl <<
    "if sys.argv[argi] == '-u':" << endl <<
    "  url = urlparse(sys.argv[argi+1])" << endl <<
    "  parts = url[1].split(':')" << endl <<
    "  host = parts[0]" << endl <<
    "  if len(parts) > 1:" << endl <<
    "    port = int(parts[1])" << endl <<
    "  else:" << endl <<
    "    port = 80" << endl <<
    "  uri = url[2]" << endl <<
    "  if url[4]:" << endl <<
    "    uri += '?%s' % url[4]" << endl <<
    "  http = True" << endl <<
    "  argi += 2" << endl <<
    endl <<
    "if sys.argv[argi] == '-f' or sys.argv[argi] == '-framed':" << endl <<
    "  framed = True" << endl <<
    "  argi += 1" << endl <<
    endl <<
    "cmd = sys.argv[argi]" << endl <<
    "args = sys.argv[argi+1:]" << endl <<
    endl <<
    "if http:" << endl <<
    "  transport = THttpClient.THttpClient(host, port, uri)" << endl <<
    "else:" << endl <<
    "  socket = TSocket.TSocket(host, port)" << endl <<
    "  if framed:" << endl <<
    "    transport = TTransport.TFramedTransport(socket)" << endl <<
    "  else:" << endl <<
    "    transport = TTransport.TBufferedTransport(socket)" << endl <<
    "protocol = TBinaryProtocol.TBinaryProtocol(transport)" << endl <<
    "client = " << service_name_ << ".Client(protocol)" << endl <<
    "transport.open()" << endl <<
    endl;

  // Generate the dispatch methods
  bool first = true;

  for (f_iter = functions.begin(); f_iter != functions.end(); ++f_iter) {
    if (first) {
      first = false;
    } else {
      f_remote << "el";
    }

    t_struct* arg_struct = (*f_iter)->get_arglist();
    const std::vector<t_field*>& args = arg_struct->get_members();
    vector<t_field*>::const_iterator a_iter;
    int num_args = args.size();

    f_remote <<
      "if cmd == '" << (*f_iter)->get_name() << "':" << endl <<
      "  if len(args) != " << num_args << ":" << endl <<
      "    print '" << (*f_iter)->get_name() << " requires " << num_args << " args'" << endl <<
      "    sys.exit(1)" << endl <<
      "  pp.pprint(client." << (*f_iter)->get_name() << "(";
    for (int i = 0; i < num_args; ++i) {
      if (args[i]->get_type()->is_string()) {
        f_remote << "args[" << i << "],";
      } else {
        f_remote << "eval(args[" << i << "]),";
      }
    }
    f_remote << "))" << endl;

    f_remote << endl;
  }
  f_remote << "else:" << endl;
  f_remote << "  print 'Unrecognized method %s' % cmd" << endl;
  f_remote << "  sys.exit(1)" << endl;
  f_remote << endl;

  f_remote << "transport.close()" << endl;

  // Close service file
  f_remote.close();

  // Make file executable, love that bitwise OR action
  chmod(f_remote_name.c_str(),
          S_IRUSR
        | S_IWUSR
        | S_IXUSR
#ifndef MINGW
        | S_IRGRP
        | S_IXGRP
        | S_IROTH
        | S_IXOTH
#endif
  );
}

/**
 * Generates a service server definition.
 *
 * @param tservice The service to generate a server for.
 */
void t_go_generator::generate_service_server(t_service* tservice) {
  // Generate the dispatch methods
  vector<t_function*> functions = tservice->get_functions();
  vector<t_function*>::iterator f_iter;

  string extends = "";
  string extends_package = "";
  string extends_processor = "";
  if (tservice->get_extends() != NULL) {
    extends_package = tservice->get_extends()->get_program()->get_name();
    extends = type_name(tservice->get_extends());
    extends_processor = extends + "Processor";
  }

  string service_name = tservice->get_name();
  string processor_name = service_name + "Processor";
  // Generate the header portion
  f_service_ <<
    "type " << processor_name << " struct {" << endl;

  indent_up();
  indent(f_service_) << "handler " << service_name << endl;
  if (!extends.empty()) {
    indent(f_service_) << "parentProcessor *" << extends_package << "." << extends_processor << endl;
  }
  indent_down();
  indent(f_service_) << "}" << endl;

  indent(f_service_) <<
    "func New" << processor_name << "(handler " << service_name << ") *" << processor_name << " {" << endl;
  indent_up();
  f_service_ <<
    indent() << "p := &" << processor_name << "{handler: handler";
  if (!extends.empty()) {
    f_service_ << ", parentProcessor: " << extends_package << ".New" << extends_processor << "(handler)";
  }
  f_service_ << "}" << endl <<
    indent() << "return p" << endl;
  indent_down();
  indent(f_service_) << "}" << endl;
  
  // Generate the server implementation
  indent(f_service_) <<
    "func (p *" << processor_name << ") Process(iprot, oprot thrift.TProtocol) (bool, *thrift.TException) {" << endl;
  indent_up();

  f_service_ <<
    indent() << "name, ttype, seqid := iprot.ReadMessageBegin()" << endl <<
    indent() << "return p.ProcessMessage(name, ttype, seqid, iprot, oprot)" << endl;
  indent_down();
  indent(f_service_) << "}" << endl;
  
  indent(f_service_) <<
    "func (p *" << processor_name << ") ProcessMessage(name string, ttype thrift.TType, seqid int32, iprot, oprot thrift.TProtocol) (bool, *thrift.TException) {" << endl;
  indent_up();

  // Switch on all the available functions
  indent(f_service_) << "switch name {" << endl;
  
  for (f_iter = functions.begin(); f_iter != functions.end(); ++f_iter) {
    string funname = (*f_iter)->get_name();
    indent(f_service_) << "case \"" << funname << "\":" << endl;
    indent_up();
    f_service_ <<
      indent() << "p.process_" << capitalize(funname) << "(seqid, iprot, oprot)" << endl;
    indent_down();
  }
  // default case for not finding a function
  // FIXME: deal with extends
  f_service_ <<
    indent() << "default:" << endl;
  indent_up();
  
  if (extends.empty()) {
    f_service_ <<
      indent() << "thrift.SkipType(iprot, thrift.TTYPE_STRUCT)" << endl <<
      indent() << "iprot.ReadMessageEnd()" << endl <<
      indent() << "err := thrift.NewTApplicationException(thrift.TAPPLICATION_EXCEPTION_UNKNOWN_METHOD, fmt.Sprintf(\"Unknown function %s\", name))" << endl <<
  		indent() << "oprot.WriteMessageBegin(name, thrift.TMESSAGETYPE_EXCEPTION, seqid)" << endl <<
  		indent() << "err.Write(oprot)" << endl <<
  		indent() << "oprot.WriteMessageEnd()" << endl <<
  		indent() << "oprot.GetTransport().Flush()" << endl <<
      indent() << "return false, err" << endl;    
  } else {
    f_service_ <<
      indent() << "return parentProcessor.ProcessMessage(name, ttype, seqid)" << endl;
  }
  indent_down();
    
  indent(f_service_) << "}" << endl;
  
  // Read end of args field, the T_STOP, and the struct close
  f_service_ <<
    indent() << "return true, nil" << endl;

  indent_down();
  f_service_ << "}" << endl;

  // Generate the process subfunctions
  for (f_iter = functions.begin(); f_iter != functions.end(); ++f_iter) {
    generate_process_function(tservice, *f_iter);
  }

  f_service_ << endl;
}

/**
 * Generates a process function definition.
 *
 * @param tfunction The function to write a dispatcher for
 */
void t_go_generator::generate_process_function(t_service* tservice,
                                               t_function* tfunction) {
  string service_name = tservice->get_name();
  string processor_name = service_name + "Processor";
  
  // Open function
  indent(f_service_) <<
    "func (p *" << processor_name << ") process_" << capitalize(tfunction->get_name()) <<
    "(seqid int32, iprot, oprot thrift.TProtocol) {" << endl;
  indent_up();

  string argsname = capitalize(tfunction->get_name()) + "_args";
  string resultname = capitalize(tfunction->get_name()) + "_result";

  f_service_ <<
    indent() << "args := New" << argsname << "()" << endl <<
    indent() << "args.Read(iprot)" << endl <<
    indent() << "iprot.ReadMessageEnd()" << endl;

  t_struct* xs = tfunction->get_xceptions();
  const std::vector<t_field*>& xceptions = xs->get_members();
  vector<t_field*>::const_iterator x_iter;

  // Declare result for non oneway function
  if (!tfunction->is_oneway()) {
    f_service_ <<
      indent() << "result := New" << resultname << "()" << endl;
  }

  // // FIXME: deal with exceptions
  // // Try block for a function with exceptions
  // if (xceptions.size() > 0) {
  //   f_service_ <<
  //     indent() << "try:" << endl;
  //   indent_up();
  // }

  // Generate the function call
  t_struct* arg_struct = tfunction->get_arglist();
  const std::vector<t_field*>& fields = arg_struct->get_members();
  vector<t_field*>::const_iterator f_iter;

  f_service_ << indent();
  if (!tfunction->is_oneway() && !tfunction->get_returntype()->is_void()) {
    f_service_ << "result.Success, _ = ";
  }
  f_service_ <<
    "p.handler." << capitalize(tfunction->get_name()) << "(";
  bool first = true;
  for (f_iter = fields.begin(); f_iter != fields.end(); ++f_iter) {
    if (first) {
      first = false;
    } else {
      f_service_ << ", ";
    }
    f_service_ << "*args." << capitalize((*f_iter)->get_name());
  }
  f_service_ << ")" << endl;

  // if (!tfunction->is_oneway() && xceptions.size() > 0) {
  //   indent_down();
  //   for (x_iter = xceptions.begin(); x_iter != xceptions.end(); ++x_iter) {
  //     f_service_ <<
  //       indent() << "except " << type_name((*x_iter)->get_type()) << ", " << (*x_iter)->get_name() << ":" << endl;
  //     if (!tfunction->is_oneway()) {
  //       indent_up();
  //       f_service_ <<
  //         indent() << "result." << (*x_iter)->get_name() << " = " << (*x_iter)->get_name() << endl;
  //       indent_down();
  //     } else {
  //       f_service_ <<
  //         indent() << "pass" << endl;
  //     }
  //   }
  // }

  // FIXME: deal w/ oneway
  // Shortcut out here for oneway functions
  if (tfunction->is_oneway()) {
    f_service_ <<
      indent() << "return" << endl;
    indent_down();
    f_service_ << "}" << endl;
    return;
  }

  f_service_ <<
    indent() << "oprot.WriteMessageBegin(\"" << tfunction->get_name() << "\", thrift.TMESSAGETYPE_REPLY, seqid)" << endl <<
    indent() << "result.Write(oprot)" << endl <<
    indent() << "oprot.WriteMessageEnd()" << endl <<
    indent() << "oprot.GetTransport().Flush()" << endl;

  // Close function
  indent_down();
  f_service_ << "}" << endl;
}

/**
 * Deserializes a field of any type.
 */
void t_go_generator::generate_deserialize_field(ofstream &out,
                                                t_field* tfield,
                                                string prefix,
                                                bool inclass) {
  t_type* type = get_true_type(tfield->get_type());

  if (type->is_void()) {
    throw "CANNOT GENERATE DESERIALIZE CODE FOR void TYPE: " +
      prefix + tfield->get_name();
  }

  string name = prefix + capitalize(tfield->get_name());

  if (type->is_struct() || type->is_xception()) {
    generate_deserialize_struct(out,
                                (t_struct*)type,
                                 name);
  } else if (type->is_container()) {
    generate_deserialize_container(out, type, name);
  } else if (type->is_base_type() || type->is_enum()) {
    indent(out) <<
      "v := iprot.";

    if (type->is_base_type()) {
      t_base_type::t_base tbase = ((t_base_type*)type)->get_base();
      switch (tbase) {
      case t_base_type::TYPE_VOID:
        throw "compiler error: cannot serialize void field in a struct: " +
          name;
        break;
      case t_base_type::TYPE_STRING:
        // FIXME: figure out if Go needs this binary thing
        if (((t_base_type*)type)->is_binary() || !gen_utf8strings_) {
          out << "ReadString()";
        } else {
          throw "compiler error: FIXME for binary TYPE_STRING";
        }
        break;
      case t_base_type::TYPE_BOOL:
        out << "ReadBool()";
        break;
      case t_base_type::TYPE_BYTE:
        out << "ReadByte()";
        break;
      case t_base_type::TYPE_I16:
        out << "ReadI16()";
        break;
      case t_base_type::TYPE_I32:
        out << "ReadI32()";
        break;
      case t_base_type::TYPE_I64:
        out << "ReadI64()";
        break;
      case t_base_type::TYPE_DOUBLE:
        out << "ReadDouble()";
        break;
      default:
        throw "compiler error: no Go name for base type " + t_base_type::t_base_name(tbase);
      }
    } else if (type->is_enum()) {
      out << "ReadI32()";
    }
    out << endl;
    indent(out) << name << " = &v" << endl;
  } else {
    printf("DO NOT KNOW HOW TO DESERIALIZE FIELD '%s' TYPE '%s'\n",
           tfield->get_name().c_str(), type->get_name().c_str());
  }
}

/**
 * Generates an unserializer for a struct, calling read()
 */
void t_go_generator::generate_deserialize_struct(ofstream &out,
                                                  t_struct* tstruct,
                                                  string prefix) {
  out <<
    indent() << prefix << " = New" << capitalize(type_name(tstruct)) << "()" << endl <<
    indent() << prefix << ".Read(iprot)" << endl;
}

/**
 * Serialize a container by writing out the header followed by
 * data and then a footer.
 */
void t_go_generator::generate_deserialize_container(ofstream &out,
                                                    t_type* ttype,
                                                    string prefix) {
  string size = tmp("_size");
  string ktype = tmp("_ktype");
  string vtype = tmp("_vtype");
  string etype = tmp("_etype");

  t_field fsize(g_type_i32, size);
  t_field fktype(g_type_byte, ktype);
  t_field fvtype(g_type_byte, vtype);
  t_field fetype(g_type_byte, etype);

  // Declare variables, read header
  if (ttype->is_map()) {
    out <<
      indent() << prefix << " = {}" << endl <<
      indent() << "(" << ktype << ", " << vtype << ", " << size << " ) = iprot.readMapBegin() " << endl;
  } else if (ttype->is_set()) {
    out <<
      indent() << prefix << " = set()" << endl <<
      indent() << "(" << etype << ", " << size << ") = iprot.readSetBegin()" << endl;
  } else if (ttype->is_list()) {
    out <<
      indent() << prefix << " = []" << endl <<
      indent() << "(" << etype << ", " << size << ") = iprot.readListBegin()" << endl;
  }

  // For loop iterates over elements
  string i = tmp("_i");
  indent(out) <<
    "for " << i << " in xrange(" << size << "):" << endl;

    indent_up();

    if (ttype->is_map()) {
      generate_deserialize_map_element(out, (t_map*)ttype, prefix);
    } else if (ttype->is_set()) {
      generate_deserialize_set_element(out, (t_set*)ttype, prefix);
    } else if (ttype->is_list()) {
      generate_deserialize_list_element(out, (t_list*)ttype, prefix);
    }

    indent_down();

  // Read container end
  if (ttype->is_map()) {
    indent(out) << "iprot.readMapEnd()" << endl;
  } else if (ttype->is_set()) {
    indent(out) << "iprot.readSetEnd()" << endl;
  } else if (ttype->is_list()) {
    indent(out) << "iprot.readListEnd()" << endl;
  }
}


/**
 * Generates code to deserialize a map
 */
void t_go_generator::generate_deserialize_map_element(ofstream &out,
                                                       t_map* tmap,
                                                       string prefix) {
  string key = tmp("_key");
  string val = tmp("_val");
  t_field fkey(tmap->get_key_type(), key);
  t_field fval(tmap->get_val_type(), val);

  generate_deserialize_field(out, &fkey);
  generate_deserialize_field(out, &fval);

  indent(out) <<
    prefix << "[" << key << "] = " << val << endl;
}

/**
 * Write a set element
 */
void t_go_generator::generate_deserialize_set_element(ofstream &out,
                                                       t_set* tset,
                                                       string prefix) {
  string elem = tmp("_elem");
  t_field felem(tset->get_elem_type(), elem);

  generate_deserialize_field(out, &felem);

  indent(out) <<
    prefix << ".add(" << elem << ")" << endl;
}

/**
 * Write a list element
 */
void t_go_generator::generate_deserialize_list_element(ofstream &out,
                                                        t_list* tlist,
                                                        string prefix) {
  string elem = tmp("_elem");
  t_field felem(tlist->get_elem_type(), elem);

  generate_deserialize_field(out, &felem);

  indent(out) <<
    prefix << ".append(" << elem << ")" << endl;
}


/**
 * Serializes a field of any type.
 *
 * @param tfield The field to serialize
 * @param prefix Name to prepend to field name
 */
void t_go_generator::generate_serialize_field(ofstream &out,
                                               t_field* tfield,
                                               string prefix) {
  t_type* type = get_true_type(tfield->get_type());

  // Do nothing for void types
  if (type->is_void()) {
    throw "CANNOT GENERATE SERIALIZE CODE FOR void TYPE: " +
      prefix + tfield->get_name();
  }

  if (type->is_struct() || type->is_xception()) {
    generate_serialize_struct(out,
                              (t_struct*)type,
                              prefix + capitalize(tfield->get_name()));
  } else if (type->is_container()) {
    generate_serialize_container(out,
                                 type,
                                 prefix + capitalize(tfield->get_name()));
  } else if (type->is_base_type() || type->is_enum()) {

    string name = "*" + prefix + capitalize(tfield->get_name());

    indent(out) <<
      "oprot.";

    if (type->is_base_type()) {
      t_base_type::t_base tbase = ((t_base_type*)type)->get_base();
      switch (tbase) {
      case t_base_type::TYPE_VOID:
        throw
          "compiler error: cannot serialize void field in a struct: " + name;
        break;
      case t_base_type::TYPE_STRING:
        if (((t_base_type*)type)->is_binary() || !gen_utf8strings_) {
          out << "WriteString(" << name << ")";
        } else {
          throw "compiler error: FIXME for binary TYPE_STRING";
        }
        break;
      case t_base_type::TYPE_BOOL:
        out << "WriteBool(" << name << ")";
        break;
      case t_base_type::TYPE_BYTE:
        out << "WriteByte(" << name << ")";
        break;
      case t_base_type::TYPE_I16:
        out << "WriteI16(" << name << ")";
        break;
      case t_base_type::TYPE_I32:
        out << "WriteI32(" << name << ")";
        break;
      case t_base_type::TYPE_I64:
        out << "WriteI64(" << name << ")";
        break;
      case t_base_type::TYPE_DOUBLE:
        out << "WriteDouble(" << name << ")";
        break;
      default:
        throw "compiler error: no Go name for base type " + t_base_type::t_base_name(tbase);
      }
    } else if (type->is_enum()) {
      out << "WriteI32(" << name << ")";
    }
    out << endl;
  } else {
    printf("DO NOT KNOW HOW TO SERIALIZE FIELD '%s%s' TYPE '%s'\n",
           prefix.c_str(),
           tfield->get_name().c_str(),
           type->get_name().c_str());
  }
}

/**
 * Serializes all the members of a struct.
 *
 * @param tstruct The struct to serialize
 * @param prefix  String prefix to attach to all fields
 */
void t_go_generator::generate_serialize_struct(ofstream &out,
                                               t_struct* tstruct,
                                               string prefix) {
  indent(out) <<
    prefix << ".Write(oprot)" << endl;
}

void t_go_generator::generate_serialize_container(ofstream &out,
                                                  t_type* ttype,
                                                  string prefix) {
  if (ttype->is_map()) {
    indent(out) <<
      "oprot.writeMapBegin(" <<
      type_to_enum(((t_map*)ttype)->get_key_type()) << ", " <<
      type_to_enum(((t_map*)ttype)->get_val_type()) << ", " <<
      "len(" << prefix << "))" << endl;
  } else if (ttype->is_set()) {
    indent(out) <<
      "oprot.writeSetBegin(" <<
      type_to_enum(((t_set*)ttype)->get_elem_type()) << ", " <<
      "len(" << prefix << "))" << endl;
  } else if (ttype->is_list()) {
    indent(out) <<
      "oprot.writeListBegin(" <<
      type_to_enum(((t_list*)ttype)->get_elem_type()) << ", " <<
      "len(" << prefix << "))" << endl;
  }

  if (ttype->is_map()) {
    string kiter = tmp("kiter");
    string viter = tmp("viter");
    indent(out) <<
      "for " << kiter << "," << viter << " in " << prefix << ".items():" << endl;
    indent_up();
    generate_serialize_map_element(out, (t_map*)ttype, kiter, viter);
    indent_down();
  } else if (ttype->is_set()) {
    string iter = tmp("iter");
    indent(out) <<
      "for " << iter << " in " << prefix << ":" << endl;
    indent_up();
    generate_serialize_set_element(out, (t_set*)ttype, iter);
    indent_down();
  } else if (ttype->is_list()) {
    string iter = tmp("iter");
    indent(out) <<
      "for " << iter << " in " << prefix << ":" << endl;
    indent_up();
    generate_serialize_list_element(out, (t_list*)ttype, iter);
    indent_down();
  }

  if (ttype->is_map()) {
    indent(out) <<
      "oprot.writeMapEnd()" << endl;
  } else if (ttype->is_set()) {
    indent(out) <<
      "oprot.writeSetEnd()" << endl;
  } else if (ttype->is_list()) {
    indent(out) <<
      "oprot.writeListEnd()" << endl;
  }
}

/**
 * Serializes the members of a map.
 *
 */
void t_go_generator::generate_serialize_map_element(ofstream &out,
                                                     t_map* tmap,
                                                     string kiter,
                                                     string viter) {
  t_field kfield(tmap->get_key_type(), kiter);
  generate_serialize_field(out, &kfield, "");

  t_field vfield(tmap->get_val_type(), viter);
  generate_serialize_field(out, &vfield, "");
}

/**
 * Serializes the members of a set.
 */
void t_go_generator::generate_serialize_set_element(ofstream &out,
                                                     t_set* tset,
                                                     string iter) {
  t_field efield(tset->get_elem_type(), iter);
  generate_serialize_field(out, &efield, "");
}

/**
 * Serializes the members of a list.
 */
void t_go_generator::generate_serialize_list_element(ofstream &out,
                                                      t_list* tlist,
                                                      string iter) {
  t_field efield(tlist->get_elem_type(), iter);
  generate_serialize_field(out, &efield, "");
}

/**
 * Generates the docstring for a given struct.
 */
void t_go_generator::generate_python_docstring(ofstream& out,
                                               t_struct* tstruct) {
  generate_python_docstring(out, tstruct, tstruct, "Attributes");
}

/**
 * Generates the docstring for a given function.
 */
void t_go_generator::generate_python_docstring(ofstream& out,
                                               t_function* tfunction) {
  generate_python_docstring(out, tfunction, tfunction->get_arglist(), "Parameters");
}

/**
 * Generates the docstring for a struct or function.
 */
void t_go_generator::generate_python_docstring(ofstream& out,
                                               t_doc*    tdoc,
                                               t_struct* tstruct,
                                               const char* subheader) {
  bool has_doc = false;
  stringstream ss;
  if (tdoc->has_doc()) {
    has_doc = true;
    ss << tdoc->get_doc();
  }

  const vector<t_field*>& fields = tstruct->get_members();
  if (fields.size() > 0) {
    if (has_doc) {
      ss << endl;
    }
    has_doc = true;
    ss << subheader << ":\n";
    vector<t_field*>::const_iterator p_iter;
    for (p_iter = fields.begin(); p_iter != fields.end(); ++p_iter) {
      t_field* p = *p_iter;
      ss << " - " << p->get_name();
      if (p->has_doc()) {
        ss << ": " << p->get_doc();
      } else {
        ss << endl;
      }
    }
  }

  if (has_doc) {
    generate_docstring_comment(out,
      "",
      "// ", ss.str(),
      "");
  }
}

/**
 * Generates the docstring for a generic object.
 */
void t_go_generator::generate_python_docstring(ofstream& out,
                                               t_doc* tdoc) {
  if (tdoc->has_doc()) {
    generate_docstring_comment(out,
      "",
      "// ", tdoc->get_doc(),
      "");
  }
}

/**
 * Declares an argument, which may include initialization as necessary.
 *
 * @param tfield The field
 */
string t_go_generator::declare_argument(t_field* tfield) {
  std::ostringstream result;
  result << tfield->get_name() << "=";
  if (tfield->get_value() != NULL) {
    result << "thrift_spec[" <<
      tfield->get_key() << "][4]";
  } else {
    result << "None";
  }
  return result.str();
}

/**
 * Renders a field default value, returns None otherwise.
 *
 * @param tfield The field
 */
string t_go_generator::render_field_default_value(t_field* tfield) {
  t_type* type = get_true_type(tfield->get_type());
  if (tfield->get_value() != NULL) {
    return render_const_value(type, tfield->get_value());
  } else {
    return "None";
  }
}

/**
 * Renders a function signature of the form 'type name(args)'
 *
 * @param tfunction Function definition
 * @return String of rendered function definition
 */
string t_go_generator::function_signature(t_function* tfunction) {
  string signature = capitalize(tfunction->get_name()) + "(" +
    argument_list(tfunction->get_arglist()) + ") ";
  
  if (!tfunction->get_returntype()->is_void()) {
    signature += "(*" + type_name(tfunction->get_returntype()) + ", *thrift.TException)";
  } else {
    signature += "*thrift.TException";
  }
  
  return signature;
}

/**
 * Renders an interface function signature of the form 'type name(args)'
 *
 * @param tfunction Function definition
 * @return String of rendered function definition
 */
string t_go_generator::function_signature_if(t_function* tfunction) {
  return function_signature(tfunction);
}


/**
 * Renders a field list
 */
string t_go_generator::argument_list(t_struct* tstruct) {
  string result = "";

  const vector<t_field*>& fields = tstruct->get_members();
  vector<t_field*>::const_iterator f_iter;
  bool first = true;
  for (f_iter = fields.begin(); f_iter != fields.end(); ++f_iter) {
    if (first) {
      first = false;
    } else {
      result += ", ";
    }
    result += (*f_iter)->get_name() + " " + type_name((*f_iter)->get_type());
  }
  return result;
}

/**
 * Converts the parse type to a Python tyoe
 */
string t_go_generator::type_to_enum(t_type* type) {
  type = get_true_type(type);

  if (type->is_base_type()) {
    t_base_type::t_base tbase = ((t_base_type*)type)->get_base();
    switch (tbase) {
    case t_base_type::TYPE_VOID:
      throw "NO T_VOID CONSTRUCT";
    case t_base_type::TYPE_STRING:
      return "thrift.TTYPE_STRING";
    case t_base_type::TYPE_BOOL:
      return "thrift.TTYPE_BOOL";
    case t_base_type::TYPE_BYTE:
      return "thrift.TTYPE_BYTE";
    case t_base_type::TYPE_I16:
      return "thrift.TTYPE_I16";
    case t_base_type::TYPE_I32:
      return "thrift.TTYPE_I32";
    case t_base_type::TYPE_I64:
      return "thrift.TTYPE_I64";
    case t_base_type::TYPE_DOUBLE:
      return "thrift.TTYPE_DOUBLE";
    }
  } else if (type->is_enum()) {
    return "thrift.TTYPE_I32";
  } else if (type->is_struct() || type->is_xception()) {
    return "thrift.TTYPE_STRUCT";
  } else if (type->is_map()) {
    return "thrift.TTYPE_MAP";
  } else if (type->is_set()) {
    return "thrift.TTYPE_SET";
  } else if (type->is_list()) {
    return "thrift.TTYPE_LIST";
  }

  throw "INVALID TYPE IN type_to_enum: " + type->get_name();
}

/** See the comment inside generate_go_struct_definition for what this is. */
string t_go_generator::type_to_spec_args(t_type* ttype) {
  while (ttype->is_typedef()) {
    ttype = ((t_typedef*)ttype)->get_type();
  }

  if (ttype->is_base_type() || ttype->is_enum()) {
    return "None";
  } else if (ttype->is_struct() || ttype->is_xception()) {
    return "(" + type_name(ttype) + ", " + type_name(ttype) + ".thrift_spec)";
  } else if (ttype->is_map()) {
    return "(" +
      type_to_enum(((t_map*)ttype)->get_key_type()) + "," +
      type_to_spec_args(((t_map*)ttype)->get_key_type()) + "," +
      type_to_enum(((t_map*)ttype)->get_val_type()) + "," +
      type_to_spec_args(((t_map*)ttype)->get_val_type()) +
      ")";

  } else if (ttype->is_set()) {
    return "(" +
      type_to_enum(((t_set*)ttype)->get_elem_type()) + "," +
      type_to_spec_args(((t_set*)ttype)->get_elem_type()) +
      ")";

  } else if (ttype->is_list()) {
    return "(" +
      type_to_enum(((t_list*)ttype)->get_elem_type()) + "," +
      type_to_spec_args(((t_list*)ttype)->get_elem_type()) +
      ")";
  }

  throw "INVALID TYPE IN type_to_spec_args: " + ttype->get_name();
}

/**
 * Maps a Thrift t_type to a C type.
 */
string t_go_generator::type_name(t_type* ttype)
{
  if (ttype->is_base_type()) {
    return base_type_name((t_base_type *) ttype);
  }

  if (ttype->is_struct() || ttype->is_xception() || ttype->is_service()) {
    return ttype->get_name();
  }

  if (ttype->is_container()) {
  }
  return "FIXME:NO_TYPE";

  // // FIXME deal with non-base_types
  // if (ttype->is_container ())
  // {
  //   string cname;
  // 
  //   t_container *tcontainer = (t_container *) ttype;
  //   if (tcontainer->has_cpp_name ())
  //   {
  //     cname = tcontainer->get_cpp_name ();
  //   } else if (ttype->is_map ()) {
  //     cname = "GHashTable *";
  //   } else if (ttype->is_set ()) {
  //     // since a set requires unique elements, use a GHashTable, and
  //     // populate the keys and values with the same data, using keys for
  //     // the actual writes and reads.
  //     // TODO: discuss whether or not to implement TSet, THashSet or GHashSet
  //     cname = "GHashTable *";
  //   } else if (ttype->is_list ()) {
  //     // TODO: investigate other implementations besides GPtrArray
  //     cname = "GPtrArray *";
  //     t_type *etype = ((t_list *) ttype)->get_elem_type ();
  //     if (etype->is_base_type ())
  //     {
  //       t_base_type::t_base tbase = ((t_base_type *) etype)->get_base ();
  //       switch (tbase)
  //       {
  //         case t_base_type::TYPE_VOID:
  //           throw "compiler error: cannot determine array type";
  //         case t_base_type::TYPE_BOOL:
  //         case t_base_type::TYPE_BYTE:
  //         case t_base_type::TYPE_I16:
  //         case t_base_type::TYPE_I32:
  //         case t_base_type::TYPE_I64:
  //         case t_base_type::TYPE_DOUBLE:
  //           cname = "GArray *";
  //           break;
  //         case t_base_type::TYPE_STRING:
  //           break;
  //         default:
  //           throw "compiler error: no array info for type";
  //       }
  //     }
  //   }
  // 
  //   if (is_const)
  //   {
  //     return "const " + cname;
  //   } else {
  //     return cname;
  //   }
  // }

  // // check for a namespace
  // string pname = this->nspace + ttype->get_name ();

  // if (is_complex_type (ttype))
  // {
  //   pname += " *";
  // }
}

/**
 * Returns the Go type that corresponds to the thrift type.
 *
 * @param tbase The base type
 * @return Explicit Go type, i.e. "int32"
 */
string t_go_generator::base_type_name(t_base_type* type) {
  t_base_type::t_base tbase = type->get_base();

  switch (tbase) {
  case t_base_type::TYPE_VOID:
    return "FIXME:void";
  case t_base_type::TYPE_STRING:
    return "string";
  case t_base_type::TYPE_BOOL:
    return "bool";
  case t_base_type::TYPE_BYTE:
    return "int8";
  case t_base_type::TYPE_I16:
    return "int16";
  case t_base_type::TYPE_I32:
    return "int32";
  case t_base_type::TYPE_I64:
    return "int64";
  case t_base_type::TYPE_DOUBLE:
    return "float64";
  default:
    throw "compiler error: no Go base type name for base type " + t_base_type::t_base_name(tbase);
  }
}


THRIFT_REGISTER_GENERATOR(go, "Go",
  "    spaces:       Use spaces instead of default tabs.\n" \
)
